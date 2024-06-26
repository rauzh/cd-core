package publish_criteria

import (
	"context"

	criteria "github.com/rauzh/cd-core/requests/criteria_controller"

	cdtime "github.com/rauzh/cd-core/time"

	"github.com/rauzh/cd-core/repo"
	"github.com/rauzh/cd-core/requests/base"
	"github.com/rauzh/cd-core/requests/publish"
)

const (
	ExplanationArtistReleaseLimit   = "More than limit releases per season"
	DiffArtistReleaseLimitPerSeason = -1
)

const ArtistReleaseLimitPerSeason criteria.CriteriaName = "No releases from artist more than limit"

var LimitPerSeason = 2

type ArtistReleaseLimitPerSeasonCriteria struct {
	publicationRepo repo.PublicationRepo
	artistRepo      repo.ArtistRepo
}

func (oarpsc *ArtistReleaseLimitPerSeasonCriteria) Name() criteria.CriteriaName {
	return ArtistReleaseLimitPerSeason
}

func (oarpsc *ArtistReleaseLimitPerSeasonCriteria) Apply(request base.IRequest) (result criteria.CriteriaDiff) {

	if err := request.Validate(publish.PubReq); err != nil {
		result.Explanation = criteria.ExplanationCantApply
		return
	}
	pubReq := request.(*publish.PublishRequest)

	ctx := context.Background()

	artist, err := oarpsc.artistRepo.GetByUserID(ctx, pubReq.ApplierID)
	if err != nil {
		result.Explanation = criteria.ExplanationCantApply
		return
	}

	pubsFromArtistLastSeason, err := oarpsc.publicationRepo.GetAllByArtistSinceDate(context.Background(),
		cdtime.RelevantPeriod(), artist.ArtistID)

	if err != nil {
		result.Explanation = criteria.ExplanationCantApply
		return
	}

	if len(pubsFromArtistLastSeason) > LimitPerSeason {
		result.Diff = DiffArtistReleaseLimitPerSeason
		result.Explanation = ExplanationArtistReleaseLimit
		return
	}

	result.Explanation = criteria.ExplanationOK

	return
}

type ArtistReleaseLimitPerSeasonCriteriaFabric struct {
	PublicationRepo repo.PublicationRepo
	ArtistRepo      repo.ArtistRepo
}

func (fabric *ArtistReleaseLimitPerSeasonCriteriaFabric) Create() (criteria.Criteria, error) {
	return &ArtistReleaseLimitPerSeasonCriteria{publicationRepo: fabric.PublicationRepo, artistRepo: fabric.ArtistRepo}, nil
}
