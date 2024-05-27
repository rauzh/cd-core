package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rauzh/cd-core/models"

	"github.com/rauzh/cd-core/repo"
)

type IArtistService interface {
	Create(*models.Artist) error
	Get(uint64) (*models.Artist, error)
	GetByUserID(id uint64) (*models.Artist, error)
	Update(*models.Artist) error
}

type ArtistService struct {
	repo   repo.ArtistRepo
	logger *slog.Logger
}

func NewArtistService(r repo.ArtistRepo, logger *slog.Logger) IArtistService {
	logger.Debug("NewArtistService")
	return &ArtistService{repo: r, logger: logger}
}

func (ars *ArtistService) Create(artist *models.Artist) error {
	if err := ars.repo.Create(context.Background(), artist); err != nil {
		ars.logger.Error("ARTIST_SERVICE Create", slog.Any("error", err))
		return fmt.Errorf("can't create artist with err %w", err)
	}
	ars.logger.Info("ARTIST_SERVICE Create", "artist_nickname", artist.Nickname)
	return nil
}

func (ars *ArtistService) Get(id uint64) (*models.Artist, error) {
	artist, err := ars.repo.Get(context.Background(), id)

	if err != nil {
		ars.logger.Error("ARTIST_SERVICE Get", slog.Any("error", err))
		return nil, fmt.Errorf("can't get artist with err %w", err)
	}
	ars.logger.Debug("ARTIST_SERVICE Get", "artist_nickname", artist.Nickname)
	return artist, nil
}

func (ars *ArtistService) GetByUserID(id uint64) (*models.Artist, error) {
	artist, err := ars.repo.GetByUserID(context.Background(), id)

	if err != nil {
		ars.logger.Error("ARTIST_SERVICE GetByUserID", slog.Any("error", err))
		return nil, fmt.Errorf("can't get artist with err %w", err)
	}
	ars.logger.Debug("ARTIST_SERVICE GetByUserID", "artist_nickname", artist.Nickname)
	return artist, nil
}

func (ars *ArtistService) Update(artist *models.Artist) error {
	if err := ars.repo.Update(context.Background(), artist); err != nil {
		ars.logger.Error("ARTIST_SERVICE Update", slog.Any("error", err), "artist_nickname", artist.Nickname)
		return fmt.Errorf("can't update artist with err %w", err)
	}
	ars.logger.Info("ARTIST_SERVICE Update", "artist_nickname", artist.Nickname)
	return nil
}
