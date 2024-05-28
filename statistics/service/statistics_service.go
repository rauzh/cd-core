package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	cdtime "github.com/rauzh/cd-core/time"

	"github.com/rauzh/cd-core/models"

	releaseService "github.com/rauzh/cd-core/release/service"
	"github.com/rauzh/cd-core/repo"
	"github.com/rauzh/cd-core/statistics/fetcher"
	ts "github.com/rauzh/cd-core/track/service"
)

type IStatisticsService interface {
	Create(*models.Statistics) error
	FetchByRelease(release *models.Release) error
	GetForTrack(uint64) ([]models.Statistics, error)
	GetByID(uint64) (*models.Statistics, error)
	GetRelevantGenre() (string, error)
	GetLatestStatForTrack(trackID uint64) (*models.Statistics, error)
}

type StatisticsService struct {
	trackService   ts.ITrackService
	releaseService releaseService.IReleaseService
	fetcher        fetcher.StatFetcher
	repo           repo.StatisticsRepo

	logger *slog.Logger
}

func NewStatisticsService(
	ts ts.ITrackService,
	f fetcher.StatFetcher,
	r repo.StatisticsRepo,
	rls releaseService.IReleaseService,
	logger *slog.Logger) IStatisticsService {
	return &StatisticsService{
		trackService:   ts,
		releaseService: rls,
		fetcher:        f,
		repo:           r,
		logger:         logger,
	}
}

func (statSvc *StatisticsService) Create(stat *models.Statistics) error {

	stat.Date = cdtime.GetToday()

	if err := statSvc.repo.Create(context.Background(), stat); err != nil {
		statSvc.logger.Error("STAT_SERVICE Create", slog.Any("error", err))
		return fmt.Errorf("can't create stats with err %w", err)
	}

	statSvc.logger.Debug("STAT_SERVICE Create", "stat_id", stat.StatID)
	return nil
}

func (statSvc *StatisticsService) GetByID(statID uint64) (*models.Statistics, error) {
	stat, err := statSvc.repo.GetByID(context.Background(), statID)

	if err != nil {
		statSvc.logger.Error("STAT_SERVICE GetByID", "stat_id", statID, slog.Any("error", err))
		return nil, fmt.Errorf("can't get stats with err %w", err)
	}
	statSvc.logger.Debug("STAT_SERVICE GetByID", "stat_id", statID)
	return stat, nil
}

func (statSvc *StatisticsService) GetForTrack(trackID uint64) ([]models.Statistics, error) {

	statSvc.logger.Info("STATS SERVICE",
		"get for track", trackID)

	stats, err := statSvc.repo.GetForTrack(context.Background(), trackID)

	if err != nil {
		statSvc.logger.Error("STAT_SERVICE GetForTrack", "track_id", trackID, slog.Any("error", err))
		return nil, fmt.Errorf("can't get stats for track %d with err %w", trackID, err)
	}

	statSvc.logger.Debug("STAT_SERVICE GetForTrack", "track_id", trackID)
	return stats, nil
}

func (statSvc *StatisticsService) FetchByRelease(release *models.Release) error {

	tracks, err := statSvc.releaseService.GetAllTracks(release)
	if err != nil {
		statSvc.logger.Error("STAT_SERVICE FetchByRelease can't get tracks", "release", release.ReleaseID, slog.Any("error", err))
		return fmt.Errorf("can't fetch stats with err %w", err)
	}

	stats, err := statSvc.fetcher.Fetch(tracks)
	if err != nil {
		statSvc.logger.Error("STAT_SERVICE FetchByRelease can't fetch", "release", release.ReleaseID, slog.Any("error", err))
		return fmt.Errorf("can't fetch stats with err %w", err)
	}

	statSvc.logger.Debug("FetchByRelease", "stats", stats)

	if len(stats) < 1 {
		statSvc.logger.Info("STAT_SERVICE FetchByRelease no stats to fetch", "release", release.ReleaseID)
		return errors.New("no stats to fetch")
	}

	for _, stat := range stats {
		stat.Date = cdtime.GetToday()
	}

	if err = statSvc.repo.CreateMany(context.Background(), stats); err != nil {
		statSvc.logger.Error("STAT_SERVICE FetchByRelease can't create many", "release_id", release.ReleaseID, slog.Any("error", err))
		return fmt.Errorf("can't create stats with err %w", err)
	}

	statSvc.logger.Info("STAT_SERVICE FetchByRelease", "release_id", release.ReleaseID)

	return nil
}

func (statSvc *StatisticsService) GetRelevantGenre() (string, error) {

	stats, err := statSvc.repo.GetAllGroupByTracksSince(context.Background(), cdtime.RelevantPeriod())
	if err != nil {
		return "", fmt.Errorf("can't get stats with err %w", err)
	}

	genres := make(map[string]uint64)
	for trackID, statsPerTrack := range *stats {
		track, err := statSvc.trackService.Get(trackID)
		if err != nil {
			return "", fmt.Errorf("can't get track %d with err %w", trackID, err)
		}

		for _, stat := range statsPerTrack {
			genres[track.Genre] += stat.Streams
		}
	}

	var relevantGenre string
	var maxStreamsPerGenre uint64

	for genre, streams := range genres {
		if streams > maxStreamsPerGenre {
			maxStreamsPerGenre = streams
			relevantGenre = genre
		}
	}

	return relevantGenre, err
}

func (statSvc *StatisticsService) GetLatestStatForTrack(trackID uint64) (*models.Statistics, error) {

	statSvc.logger.Info("STATS SERVICE",
		"get latest stat for track", trackID)

	stats, err := statSvc.GetForTrack(trackID)
	if err != nil {
		return nil, err
	}

	statSvc.logger.Debug("STATS SERVICE",
		"get latest stat for track", trackID,
		"stats len", len(stats))

	latestStatDate := stats[0].Date
	latestStat := stats[0]

	for _, stat := range stats {
		if stat.Date.After(latestStatDate) {
			latestStatDate = stat.Date
			latestStat = stat
		}
	}

	return &latestStat, nil
}
