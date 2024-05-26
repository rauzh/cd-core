package sign_contract

import (
	"context"

	"github.com/rauzh/cd-core/requests/base"
	"github.com/rauzh/cd-core/requests/sign_contract"
	"github.com/rauzh/cd-core/requests/sign_contract/errors"
)

func (handler *SignContractProceedToManagerHandler) proceedToManager(signReq *sign_contract.SignContractRequest) error {
	signReq.Status = base.OnApprovalRequest

	ctx := context.Background()

	managerID, err := handler.mngRepo.GetRandManagerID(ctx)
	if err != nil {
		return errors.ErrCantFindManager
	}

	//fmt.Println("!!! MNG ID", managerID)

	signReq.ManagerID = managerID

	return handler.signReqRepo.Update(ctx, signReq)
}
