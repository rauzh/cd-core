package repo

import (
	"context"
	"time"

	"github.com/rauzh/cd-core/models"
)

//go:generate mockery --name StatisticsRepo --with-expecter
type StatisticsRepo interface {
	Create(context.Context, *models.Statistics) error
	GetForTrack(context.Context, uint64) ([]models.Statistics, error)
	GetByID(context.Context, uint64) (*models.Statistics, error)
	GetAllGroupByTracksSince(ctx context.Context, date time.Time) (*map[uint64][]models.Statistics, error)
	CreateMany(context.Context, []models.Statistics) error
}
