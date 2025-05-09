package service

import (
	"context"
	"errors"
	"github.com/SwanHtetAungPhyo/wolftagon/internal/model"
	"github.com/SwanHtetAungPhyo/wolftagon/pkg/jwt_provider"
	"github.com/redis/go-redis/v9"
	"github.com/resend/resend-go/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	accessTokenBlacklistPrefix  = "black_access:"
	refreshTokenBlacklistPrefix = "black_refresh:"
)

var _UserServiceBehaviour = (*UserService)(nil)

func (u *UserService) Login(req *model.LoginUserRequest) (*model.LoginResponse, *string, error) {
	u.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("Attempting user login")

	user, err := u.repo.GetByEmail(req.Email)
	if err != nil {
		u.log.WithError(err).Warn("Login failed - user not found")
		return nil, nil, errors.New("invalid credentials")
	}

	//if !user.Verified {
	//	u.log.WithFields(logrus.Fields{
	//		"email":   req.Email,
	//		"message": "Login failed - email not verified",
	//	})
	//	return nil, nil, errors.New("you need to verify your email first")
	//}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		u.log.WithError(err).Warn("Login failed - invalid password")
		return nil, nil, errors.New("invalid credentials")
	}

	accessToken, err := jwt_provider.JwtTokenGenerator(0, user.UserID.String(), user.Role.RoleName)
	if err != nil {
		u.log.WithError(err).Error("Failed to generate access token")
		return nil, nil, err
	}

	refreshToken, err := jwt_provider.JwtTokenGenerator(1, user.UserID.String(), user.Role.RoleName)
	if err != nil {
		u.log.WithError(err).Error("Failed to generate refresh token")
		return nil, nil, err
	}

	u.log.WithFields(logrus.Fields{
		"user_id": user.UserID,
		"email":   user.Email,
	}).Info("User logged in successfully")

	return &model.LoginResponse{
		Message: "Login successful",
		Token:   accessToken,
		EmbeddedUserDataInLoginSuccess: model.EmbeddedUserDataInLoginSuccess{
			UserId:    user.UserID,
			Email:     user.Email,
			FirstName: user.FirstName,
			RoleName:  user.Role.RoleName,
		},
	}, &refreshToken, nil
}

func (u *UserService) Register(req *model.RegisterUserRequest) (*model.RegisterResponse, error) {
	u.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("attempting user registration")
	existingUser, err := u.repo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("database error")
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.WithError(err).Error("failed to hash password")
		return nil, errors.New("password processing error")
	}

	user := &model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Age:       req.Age,
		Verified:  false,
	}

	verificationToken := u.EmailVerificationToken()

	err = u.repo.Create(user, req.RoleName)
	if err != nil {
		return nil, err
	}

	go func() {
		u.sendVerificationEmail(verificationToken, user.Email)

	}()

	if err := u.redis.Set(context.Background(),
		user.Email,
		verificationToken,
		5*time.Minute).Err(); err != nil {
		u.log.WithError(err).Error("failed to store verification token")
		return nil, errors.New("failed to initiate verification process")
	}

	return &model.RegisterResponse{
		Message: "Registration successful. Verification email sent.",
	}, nil
}
func (u *UserService) generateAndStoreVerificationToken(email string) {
	token := u.EmailVerificationToken()

	u.emailWG.Add(2)

	go func() {
		defer u.emailWG.Done()
		u.sendVerificationEmail(token, email)
	}()

	go func() {
		defer u.emailWG.Done()
		if err := u.redis.Set(context.Background(), email, token, 5*time.Minute).Err(); err != nil {
			u.log.WithError(err).Error("Failed to store verification token in Redis")
		}
	}()
}
func (u *UserService) Verify(req *model.EmailVerificationReq) (bool, error) {
	u.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("attempting email verification")

	tokenInRedis, err := u.redis.Get(context.Background(), req.Email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			u.log.WithFields(logrus.Fields{
				"email": req.Email,
			}).Warn("verification token expired or not found")
			return false, errors.New("verification token expired or invalid")
		}
		u.log.WithError(err).Error("redis error fetching token")
		return false, errors.New("verification service unavailable")
	}

	if !strings.EqualFold(strings.TrimSpace(req.Code), strings.TrimSpace(tokenInRedis)) {
		u.log.Warn("verification code mismatch")
		return false, errors.New("invalid verification code")
	}

	if err := u.repo.MarkAsVerified(req.Email); err != nil {
		u.log.WithError(err).Error("failed to mark user as verified")
		return false, errors.New("failed to complete verification")
	}

	if err := u.redis.Del(context.Background(), req.Email).Err(); err != nil {
		u.log.WithError(err).Warn("failed to delete used verification token")
	}

	u.log.WithField("email", req.Email).Info("email verified successfully")
	return true, nil
}

func (u *UserService) Logout(refreshToken, accessToken, userId string) {
	u.log.WithFields(logrus.Fields{
		"email": userId,
	}).Info("Processing user logout")

	u.redis.Set(context.Background(),
		"blacklist:"+userId,
		refreshToken+":"+accessToken,
		5*time.Minute,
	)
}

func (u *UserService) sendVerificationEmail(token, email string) {
	apiKey := os.Getenv("RES_API_KEY")
	if apiKey == "" {
		u.log.Error("Resend API key not configured")
		return
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{email},
		Html:    "<p>Your verification code is: <strong>" + token + "</strong></p>",
		Subject: "Email Verification Code for platform",

		ReplyTo: "replyto@example.com",
	}
	_, err := client.Emails.Send(params)
	if err != nil {
		u.log.WithFields(logrus.Fields{
			"email": email,
			"error": err,
		}).Error("Failed to send verification email")
		return
	}

	u.log.WithFields(logrus.Fields{
		"email": email,
	}).Info("Verification email sent successfully")
}

func (u *UserService) EmailVerificationToken() string {
	const charSet = "0123456789"
	const length = 6
	token := make([]byte, length)

	for i := range token {
		num := charSet[rand.Intn(len(charSet))]
		token[i] = num
	}

	return string(token)
}

func (u *UserService) WaitForEmailOperations() {
	u.emailWG.Wait()
}
func (u *UserService) BlacklistTokens(userID, accessToken, refreshToken string) error {
	ctx := context.Background()

	pipe := u.redis.Pipeline()

	pipe.Set(ctx,
		"black_access:"+userID,
		accessToken,
		1*time.Hour)

	if refreshToken != "" {
		pipe.Set(ctx,
			"black_refresh:"+userID,
			refreshToken,
			24*time.Hour)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		u.log.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("failed to execute redis pipeline for token blacklisting")
		return errors.New("failed to blacklist tokens")
	}

	return nil
}

func (u *UserService) IsTokenBlacklisted(userID, token string, tokenType int) (bool, error) {
	ctx := context.Background()
	var key string

	switch tokenType {
	case 0:
		key = "black_access:" + userID
	case 1:
		key = "black_refresh:" + userID
	default:
		return false, errors.New("invalid token type")
	}

	blacklistedToken, err := u.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return blacklistedToken == token, nil
}
