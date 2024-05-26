package fetcher

import "github.com/rauzh/cd-core/models"

//go:generate mockery --name StatFetcher --with-expecter
type StatFetcher interface {
	Fetch(tracks []models.Track) ([]models.Statistics, error)
}
