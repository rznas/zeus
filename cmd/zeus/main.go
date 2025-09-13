package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	docs "github.com/rznas/zeus/docs"
	"github.com/rznas/zeus/internal/config"
	"github.com/rznas/zeus/internal/db"
	"github.com/rznas/zeus/internal/middleware"
	"github.com/rznas/zeus/internal/models"
	"github.com/rznas/zeus/internal/repositories"
	"github.com/rznas/zeus/internal/routes"
	"github.com/rznas/zeus/internal/services"
)

// @title Zeus API
// @version 1.0
// @description Fiber + GORM + Redis user service with OTP and JWT.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	// Set swagger info
	docs.SwaggerInfo.Host = "localhost:" + cfg.App.Port
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "Zeus API"
	docs.SwaggerInfo.Description = "Fiber + GORM + Redis user service with OTP and JWT."

	// Setup Postgres
	gormDB, err := db.NewPostgres(db.PostgresOptions{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		DB:       cfg.Postgres.DB,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
	})
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}
	// Migrate
	if err := gormDB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Setup Redis
	redisClient := db.NewRedisClient(db.RedisOptions{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	if err := db.RedisPing(context.Background(), redisClient); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	// Services
	otpSvc := services.NewOTPService(redisClient, cfg.App.OTPTTLSeconds, cfg.App.OTPRatePerMin, cfg.App.OTPRateLimitSeconds)
	jwtSvc := services.NewJWTService(cfg.App.JWTSecret, cfg.App.JWTExpiresMinutes)

	// repository
	userRepo := repositories.NewUserRepository(gormDB)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
		MaxAge:           int((12 * time.Hour).Seconds()),
	}))
	app.Use(limiter.New(limiter.Config{Max: cfg.App.RateLimitPerMin}))

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health
	// @Summary Health check
	// @Tags Health
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Routes
	auth := &routes.AuthHandlers{DB: gormDB, OTP: otpSvc, JWT: jwtSvc, Env: cfg.App.Env}
	users := &routes.UsersHandlers{UserRepo: userRepo}

	api := app.Group("/api")
	auth.RegisterRoutes(api.Group("/auth"))
	// Protected group
	protected := api.Group("", middleware.AuthMiddleware(jwtSvc))
	users.RegisterRoutes(protected)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("starting server on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
