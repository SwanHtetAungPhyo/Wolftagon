package jwt_provider

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func JwtTokenGenerator(tokenType int, userId, userRole string) (string, error) {
	now := time.Now()
	var expiration int64
	switch tokenType {
	case AccessToken:
		expiration = now.Add(time.Minute * 30).Unix()
	case RefreshToken:
		expiration = now.Add(time.Hour * 72).Unix()
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtClaims := jwt.MapClaims{
		"sub":  userId,
		"role": userRole,
		"type": tokenType,
		"iat":  now.Unix(),
		"exp":  expiration,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
