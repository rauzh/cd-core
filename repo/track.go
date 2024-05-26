package repo

import (
	"context"

	"github.com/rauzh/cd-core/models"
)

//go:generate mockery --name TrackRepo --with-expecter
type TrackRepo interface {
	Create(context.Context, *models.Track) (uint64, error)
	Get(context.Context, uint64) (*models.Track, error)
}
