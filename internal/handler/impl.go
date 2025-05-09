package handler

import (
	"github.com/SwanHtetAungPhyo/wolftagon/internal/model"
	"github.com/SwanHtetAungPhyo/wolftagon/pkg/jwt_provider"
	"github.com/SwanHtetAungPhyo/wolftagon/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var validate = validator.New()

func (u UserHandler) Login(c *fiber.Ctx) error {
	var req model.LoginUserRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format", err)
	}

	if err := validate.Struct(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	loginResponse, refreshToken, err := u.srv.Login(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication failed", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", loginResponse)
}

func (u UserHandler) Register(c *fiber.Ctx) error {
	var req model.RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format", err)
	}

	if err := validate.Struct(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	response, err := u.srv.Register(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Registration failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Registration successful", response)
}

func (u UserHandler) Verify(c *fiber.Ctx) error {
	var req model.EmailVerificationReq

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format", err)
	}

	if err := validate.Struct(req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	verified, err := u.srv.Verify(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Verification failed", err)
	}

	if !verified {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid verification code", nil)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Email verified successfully", nil)
}

func (u UserHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	accessToken := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	userID := c.Locals("userId").(string)
	userRole := c.Locals("role").(string)

	if userID == "" {
		u.log.WithField("location", "refresh_handler").Warn("missing user ID in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", nil)
	}

	if refreshToken == "" {
		u.log.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   "missing_refresh_token",
		}).Warn("refresh attempt with missing token")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Refresh token required", nil)
	}

	isBlacklisted, err := u.srv.IsTokenBlacklisted(userID, refreshToken, 1)
	if err != nil {
		u.log.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("failed to check token blacklist")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Token verification failed", nil)
	}

	if isBlacklisted {
		u.log.WithFields(logrus.Fields{
			"user_id": userID,
		}).Warn("attempt to use blacklisted refresh token")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token", nil)
	}

	if err := u.srv.BlacklistTokens(userID, accessToken, ""); err != nil {
		u.log.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("failed to blacklist old access token")
	}

	newAccessToken, err := jwt_provider.JwtTokenGenerator(0, userID, userRole)
	if err != nil {
		u.log.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("failed to generate new access token")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to refresh token", nil)
	}

	newRefreshToken, err := jwt_provider.JwtTokenGenerator(1, userID, userRole)

	u.log.WithField("user_id", userID).Info("token refreshed successfully")
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	return utils.SuccessResponse(c, fiber.StatusOK, "Token refreshed", fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
func (u UserHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	accessToken := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	userID := c.Locals("userId").(string)

	if userID == "" {
		u.log.WithField("location", "logout_handler").Warn("missing user ID in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", nil)
	}

	if refreshToken == "" || accessToken == "" {
		u.log.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   "missing_tokens",
		}).Warn("logout attempt with missing tokens")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Authentication tokens required", nil)
	}

	if err := u.srv.BlacklistTokens(userID, accessToken, refreshToken); err != nil {
		u.log.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Error("failed to blacklist tokens during logout")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Logout processing failed", nil)
	}

	c.ClearCookie("refresh_token")

	u.log.WithField("user_id", userID).Info("user logged out successfully")
	return utils.SuccessResponse(c, fiber.StatusOK, "Logout successful", nil)
}
