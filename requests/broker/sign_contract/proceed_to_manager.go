package sign_contract

import (
	"context"
	"log/slog"

	"github.com/rauzh/cd-core/requests/base"
	"github.com/rauzh/cd-core/requests/sign_contract"
	"github.com/rauzh/cd-core/requests/sign_contract/errors"
)

func (handler *SignContractProceedToManagerHandler) proceedToManager(signReq *sign_contract.SignContractRequest) error {
	signReq.Status = base.OnApprovalRequest

	ctx := context.Background()

	managerID, err := handler.mngRepo.GetRandManagerID(ctx)
	if err != nil {
		handler.logger.Error("SIGN_HANDLER proceedToManager", slog.Any("error", err))
		return errors.ErrCantFindManager
	}

	signReq.ManagerID = managerID

	err = handler.signReqRepo.Update(ctx, signReq)
	if err != nil {
		handler.logger.Error("SIGN_HANDLER proceedToManager", slog.Any("error", err))
		return err
	}
	handler.logger.Info("SIGN_HANDLER proceedToManager", "signreq_manager", signReq.ManagerID)
	return nil
}
