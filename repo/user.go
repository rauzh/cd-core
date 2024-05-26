package repo

import (
	"context"

	"github.com/rauzh/cd-core/models"
)

//go:generate mockery --name UserRepo --with-expecter
type UserRepo interface {
	Create(context.Context, *models.User) error
	GetByEmail(context.Context, string) (*models.User, error)
	Get(context.Context, uint64) (*models.User, error)
	Update(context.Context, *models.User) error
	UpdateType(ctx context.Context, userID uint64, typ models.UserType) error
	GetForAdmin(ctx context.Context) ([]models.User, error)

	SetRole(ctx context.Context, role models.UserType) error
}
