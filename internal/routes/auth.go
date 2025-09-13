package routes

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/rznas/zeus/internal/models"
	"github.com/rznas/zeus/internal/services"
)

type AuthHandlers struct {
	DB  *gorm.DB
	OTP *services.OTPService
	JWT *services.JWTService
	Env string
}

type phoneReq struct {
	Phone string `json:"phone"`
}

type otpVerifyReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (h *AuthHandlers) RegisterRoutes(r fiber.Router) {
	// Merge login with OTP request
	r.Post("/login", h.requestOTP)
	r.Post("/otp/verify", h.verifyOTP)
}

func normalizePhone(p string) string { return strings.TrimSpace(p) }

// requestOTP
// @Summary Login (request OTP)
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body phoneReq true "Phone"
// @Success 200 {object} map[string]any
// @Router /api/auth/login [post]
func (h *AuthHandlers) requestOTP(c *fiber.Ctx) error {
	var req phoneReq
	if err := c.BodyParser(&req); err != nil || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "phone required"})
	}
	phone := normalizePhone(req.Phone)
	code, err := h.OTP.Generate(c.Context(), phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "otp error"})
	}
	if h.Env == "development" {
		log.Printf("DEV OTP for %s: %s", phone, code)
	}
	return c.JSON(fiber.Map{"sent": true})
}

// verifyOTP
// @Summary Verify OTP (register/login)
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body otpVerifyReq true "Verify"
// @Success 200 {object} map[string]string
// @Router /api/auth/otp/verify [post]
func (h *AuthHandlers) verifyOTP(c *fiber.Ctx) error {
	var req otpVerifyReq
	if err := c.BodyParser(&req); err != nil || req.Phone == "" || req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "phone and code required"})
	}
	phone := normalizePhone(req.Phone)
	ok, err := h.OTP.Verify(c.Context(), phone, req.Code)
	if err != nil || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid code"})
	}
	// Create user if not exists
	u := models.User{Phone: phone}
	if err := h.DB.WithContext(context.Background()).FirstOrCreate(&u, models.User{Phone: phone}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	tok, err := h.JWT.Generate(u.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "jwt error"})
	}
	return c.JSON(fiber.Map{"token": tok})
}
