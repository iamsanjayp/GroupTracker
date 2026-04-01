package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"grouptracker/internal/config"
)

type JWTClaims struct {
	UserID uint64 `json:"sub"`
	TeamID uint64 `json:"team_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(cfg *config.Config, userID, teamID uint64, role string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		TeamID: teamID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTAccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "grouptracker",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func GenerateRefreshToken(cfg *config.Config, userID uint64) (string, time.Time, error) {
	expiry := time.Now().Add(cfg.JWTRefreshExpiry)
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(expiry),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "grouptracker-refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	return signed, expiry, err
}

func ParseRefreshToken(cfg *config.Config, tokenStr string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}
	return claims, nil
}

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		tokenStr := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("teamID", claims.TeamID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) uint64 {
	if v, ok := c.Locals("userID").(uint64); ok {
		return v
	}
	return 0
}

func GetTeamID(c *fiber.Ctx) uint64 {
	if v, ok := c.Locals("teamID").(uint64); ok {
		return v
	}
	return 0
}

func GetRole(c *fiber.Ctx) string {
	if v, ok := c.Locals("role").(string); ok {
		return v
	}
	return ""
}
