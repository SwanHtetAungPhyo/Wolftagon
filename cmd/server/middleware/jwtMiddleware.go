package middleware

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

const (
	accessTokenBlacklistPrefix = "black_access:"
)

func RoleMiddleware(redisClient *redis.Client, logger *logrus.Logger, allowedRoles ...string) fiber.Handler {
	roleSet := make(map[string]bool)
	for _, r := range allowedRoles {
		roleSet[r] = true
	}

	return func(c *fiber.Ctx) error {
		tokenStr, userID, err := extractAndValidateToken(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		if isBlacklisted, err := checkTokenBlacklist(redisClient, userID, tokenStr); err != nil {
			logger.WithError(err).Error("Failed to check token blacklist")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		} else if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token revoked"})
		}

		claims, err := parseTokenClaims(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		role, ok := claims["role"].(string)
		if !ok || !roleSet[role] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		c.Locals("role", role)
		c.Locals("userId", claims["sub"])
		return c.Next()
	}
}

func JwtMiddleware(redisClient *redis.Client, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr, userID, err := extractAndValidateToken(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		if isBlacklisted, err := checkTokenBlacklist(redisClient, userID, tokenStr); err != nil {
			logger.WithError(err).Error("Failed to check token blacklist")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		} else if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token revoked"})
		}

		claims, err := parseTokenClaims(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		c.Locals("role", claims["role"])
		c.Locals("userId", claims["sub"])
		return c.Next()
	}
}

func extractAndValidateToken(c *fiber.Ctx) (string, string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", "", errors.New("Missing or invalid Authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("Invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", "", errors.New("Invalid user ID in token")
	}

	return tokenStr, userID, nil
}

func checkTokenBlacklist(redisClient *redis.Client, userID, token string) (bool, error) {
	blacklistedToken, err := redisClient.Get(context.Background(), accessTokenBlacklistPrefix+userID).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return blacklistedToken == token, nil
}

func parseTokenClaims(tokenStr string) (jwt.MapClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
