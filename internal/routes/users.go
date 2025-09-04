package routes

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/reza/zeus/internal/models"
)

type UsersHandlers struct {
	DB *gorm.DB
}

func (h *UsersHandlers) RegisterRoutes(r fiber.Router) {
	r.Get("/users", h.listUsers)
}

// listUsers
// @Summary List users
// @Tags Users
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Success 200 {object} map[string]any
// @Security BearerAuth
// @Router /api/users [get]
func (h *UsersHandlers) listUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	var users []models.User
	var total int64
	q := h.DB.WithContext(context.Background()).Model(&models.User{})
	if err := q.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	if err := q.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	return c.JSON(fiber.Map{
		"data":      users,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}
