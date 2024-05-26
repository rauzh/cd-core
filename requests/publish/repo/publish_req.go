package repo

import (
	"context"

	"github.com/rauzh/cd-core/requests/publish"
)

//go:generate mockery --name PublishRequestRepo --with-expecter
type PublishRequestRepo interface {
	Create(context.Context, *publish.PublishRequest) error
	Get(ctx context.Context, id uint64) (*publish.PublishRequest, error)
	Update(context.Context, *publish.PublishRequest) error
	SetMeta(context.Context, *publish.PublishRequest) error
}
