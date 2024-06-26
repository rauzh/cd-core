package publish

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rauzh/cd-core/requests/base"
	criteria "github.com/rauzh/cd-core/requests/criteria_controller"
	"github.com/rauzh/cd-core/requests/publish"
)

func (handler *PublishProceedToManagerConsumerHandler) proceedToManager(pubReq *publish.PublishRequest) error {

	ctx := context.Background()

	pubReq.Status = base.ProcessingRequest

	err := handler.publishRepo.Update(ctx, pubReq)
	if err != nil {
		handler.logger.Error("PUBLISH_HANDLER proceedToManager", slog.Any("error", err))
		return fmt.Errorf("cant proceed publish request to manager: update repo with err %w", err)
	}

	handler.computeDegree(pubReq)

	artist, err := handler.artistRepo.GetByUserID(ctx, pubReq.ApplierID)
	if err != nil {
		handler.logger.Error("PUBLISH_HANDLER proceedToManager", slog.Any("error", err))
		return fmt.Errorf("cant proceed publish request to manager: get artist with err %w", err)
	}

	pubReq.ManagerID = artist.ManagerID
	pubReq.Status = base.OnApprovalRequest

	err = handler.publishRepo.Update(ctx, pubReq)
	if err != nil {
		handler.logger.Error("PUBLISH_HANDLER proceedToManager", slog.Any("error", err))
		return err
	}
	handler.logger.Info("PUBLISH_HANDLER proceedToManager", "pubreq_manager", pubReq.ManagerID)
	return nil
}

func (handler *PublishProceedToManagerConsumerHandler) computeDegree(pubReq *publish.PublishRequest) {

	summaryDiff := handler.criterias.Apply(pubReq)

	pubReq.Grade = summaryDiff.ResultDiff
	for criteriaName, criteriaDiff := range summaryDiff.ResultExplanation {
		pubReq.Description += criteria.DiffToString(criteriaName, criteriaDiff.Explanation, criteriaDiff.Diff)
	}
}
