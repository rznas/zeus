package repositories

import (
	"context"

	"github.com/rznas/zeus/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	List(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}
