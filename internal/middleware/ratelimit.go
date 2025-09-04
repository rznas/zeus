package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func GlobalRateLimiter(maxPerMinute int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxPerMinute,
		Expiration: time.Minute,
	})
}
