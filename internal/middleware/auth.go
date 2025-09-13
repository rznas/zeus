package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rznas/zeus/internal/services"
)

type contextKey string

const (
	ContextUserID contextKey = "user_id"
)

func AuthMiddleware(jwtSvc *services.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
		token = strings.TrimSpace(token)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}
		claims, err := jwtSvc.Parse(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		uid, err := uuid.Parse(claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid uid"})
		}
		c.Locals(string(ContextUserID), uid)
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	v := c.Locals(string(ContextUserID))
	if v == nil {
		return uuid.Nil, false
	}
	uid, ok := v.(uuid.UUID)
	return uid, ok
}
