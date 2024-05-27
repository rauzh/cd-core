package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rauzh/cd-core/models"

	"github.com/rauzh/cd-core/repo"
)

type IPublicationService interface {
	Create(*models.Publication) error
	Update(*models.Publication) error
	Get(uint64) (*models.Publication, error)
	GetAllByDate(date time.Time) ([]models.Publication, error)
	GetAllByManager(managerID uint64) ([]models.Publication, error)
	GetAllByArtistSinceDate(date time.Time, artistID uint64) ([]models.Publication, error)
}

type PublicationService struct {
	repo   repo.PublicationRepo
	logger *slog.Logger
}

func NewPublicationService(repo repo.PublicationRepo, logger *slog.Logger) IPublicationService {
	return &PublicationService{repo: repo, logger: logger}
}

func (pbcSvc *PublicationService) Create(publication *models.Publication) error {
	if err := pbcSvc.repo.Create(context.Background(), publication); err != nil {
		pbcSvc.logger.Error("PUBLICATION_SERVICE Create", slog.Any("error", err))
		return fmt.Errorf("can't create publication info with error %w", err)
	}
	pbcSvc.logger.Info("PUBLICATION_SERVICE Create", "publication_ID", publication.PublicationID)
	return nil
}

func (pbcSvc *PublicationService) Update(publication *models.Publication) error {
	if err := pbcSvc.repo.Update(context.Background(), publication); err != nil {
		pbcSvc.logger.Error("PUBLICATION_SERVICE Update", slog.Any("error", err))
		return fmt.Errorf("can't update publication info with error %w", err)
	}
	pbcSvc.logger.Info("PUBLICATION_SERVICE Update", "publication_ID", publication.PublicationID)
	return nil
}

func (pbcSvc *PublicationService) Get(id uint64) (*models.Publication, error) {
	publication, err := pbcSvc.repo.Get(context.Background(), id)

	if err != nil {
		pbcSvc.logger.Error("PUBLICATION_SERVICE Get", slog.Any("error", err))
		return nil, fmt.Errorf("can't get publication info with error %w", err)
	}
	pbcSvc.logger.Debug("PUBLICATION_SERVICE Get", "publication_ID", publication.PublicationID)
	return publication, nil
}

func (pbcSvc *PublicationService) GetAllByDate(date time.Time) ([]models.Publication, error) {
	publications, err := pbcSvc.repo.GetAllByDate(context.Background(), date)

	if err != nil {
		pbcSvc.logger.Error("PUBLICATION_SERVICE GetAllByDate", slog.Any("error", err))
		return nil, fmt.Errorf("can't get publications info with error %w", err)
	}
	pbcSvc.logger.Debug("PUBLICATION_SERVICE GetAllByDate", "publications_len", len(publications))
	return publications, nil
}

func (pbcSvc *PublicationService) GetAllByManager(mng uint64) ([]models.Publication, error) {
	publications, err := pbcSvc.repo.GetAllByManager(context.Background(), mng)
	if err != nil {
		pbcSvc.logger.Error("PUBLICATION_SERVICE GetAllByManager", slog.Any("error", err))
		return nil, fmt.Errorf("can't get publications info with error %w", err)
	}
	pbcSvc.logger.Debug("PUBLICATION_SERVICE GetAllByManager", "publications_len", len(publications))
	return publications, nil
}

func (pbcSvc *PublicationService) GetAllByArtistSinceDate(date time.Time, artistID uint64) ([]models.Publication, error) {
	publications, err := pbcSvc.repo.GetAllByArtistSinceDate(context.Background(), date, artistID)

	if err != nil {
		pbcSvc.logger.Error("PUBLICATION_SERVICE GetAllByArtistSinceDate", slog.Any("error", err))
		return nil, fmt.Errorf("can't get publications info with error %w", err)
	}
	pbcSvc.logger.Debug("PUBLICATION_SERVICE GetAllByArtistSinceDate", "publications_len", len(publications))
	return publications, nil
}
