package publish_criteria

import (
	"context"

	criteria "github.com/rauzh/cd-core/requests/criteria_controller"

	"github.com/rauzh/cd-core/repo"
	"github.com/rauzh/cd-core/requests/base"
	"github.com/rauzh/cd-core/requests/publish"
)

const (
	ReleasesPerDayLimit                         = 1
	OneReleasePerDay      criteria.CriteriaName = "No releases that day"
	ExplanationOneRelease                       = "More than one release per day"
	DiffOneRelease                              = -1
)

type OneReleasePerDayCriteria struct {
	publicationRepo repo.PublicationRepo
}

func (orpdc *OneReleasePerDayCriteria) Name() criteria.CriteriaName {
	return OneReleasePerDay
}

func (orpdc *OneReleasePerDayCriteria) Apply(request base.IRequest) (result criteria.CriteriaDiff) {

	if err := request.Validate(publish.PubReq); err != nil {
		result.Explanation = criteria.ExplanationCantApply
		return
	}
	pubReq := request.(*publish.PublishRequest)

	pubsThatDay, err := orpdc.publicationRepo.GetAllByDate(context.Background(), pubReq.ExpectedDate)
	if err != nil {
		result.Explanation = criteria.ExplanationCantApply
		return
	}

	if len(pubsThatDay) > ReleasesPerDayLimit {
		result.Diff = DiffOneRelease
		result.Explanation = ExplanationOneRelease
		return
	}

	result.Explanation = criteria.ExplanationOK

	return
}

type OneReleasePerDayCriteriaFabric struct {
	PublicationRepo repo.PublicationRepo
}

func (fabric *OneReleasePerDayCriteriaFabric) Create() (criteria.Criteria, error) {
	return &OneReleasePerDayCriteria{publicationRepo: fabric.PublicationRepo}, nil
}
