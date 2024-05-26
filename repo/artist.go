package repo

import (
	"context"

	"github.com/rauzh/cd-core/models"
)

//go:generate mockery --name ArtistRepo --with-expecter
type ArtistRepo interface {
	Create(context.Context, *models.Artist) error
	Get(context.Context, uint64) (*models.Artist, error)
	GetByUserID(context.Context, uint64) (*models.Artist, error)
	Update(context.Context, *models.Artist) error
}
