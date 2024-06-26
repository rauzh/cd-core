package usecase

import (
	"context"
	"fmt"
	"log/slog"

	cdtime "github.com/rauzh/cd-core/time"

	"github.com/rauzh/cd-core/models"

	repo "github.com/rauzh/cd-core/repo"
	"github.com/rauzh/cd-core/requests/base"
	"github.com/rauzh/cd-core/requests/broker"
	"github.com/rauzh/cd-core/requests/broker/broker_dto"
	signContractBroker "github.com/rauzh/cd-core/requests/broker/sign_contract"
	"github.com/rauzh/cd-core/requests/sign_contract"
	signContractRepo "github.com/rauzh/cd-core/requests/sign_contract/repo"
	"github.com/rauzh/cd-core/transactor"
)

type SignContractRequestUseCase struct {
	userRepo   repo.UserRepo
	artistRepo repo.ArtistRepo
	transactor transactor.Transactor
	scBroker   broker.IBroker

	repo signContractRepo.SignContractRequestRepo

	logger *slog.Logger
}

func NewSignContractRequestUseCase(
	usrRepo repo.UserRepo,
	artRepo repo.ArtistRepo,
	transactor transactor.Transactor,
	scBroker broker.IBroker,
	repo signContractRepo.SignContractRequestRepo,
	logger *slog.Logger,
) (base.IRequestUseCase, error) {

	sctUseCase := &SignContractRequestUseCase{
		userRepo:   usrRepo,
		artistRepo: artRepo,
		repo:       repo,
		transactor: transactor,
		scBroker:   scBroker,
		logger:     logger,
	}

	return sctUseCase, nil
}

func (sctUseCase *SignContractRequestUseCase) Apply(request base.IRequest) error {

	if err := request.Validate(sign_contract.SignRequest); err != nil {
		return err
	}
	signReq := request.(*sign_contract.SignContractRequest)

	base.InitDateStatus(&signReq.Request)

	if err := sctUseCase.repo.Create(context.Background(), signReq); err != nil {
		return fmt.Errorf("can't apply sign contract request with err %w", err)
	}

	if err := sctUseCase.sendProceedToManagerMSG(signReq); err != nil {
		return err
	}

	sctUseCase.logger.Info("SIGNREQ_UC Apply", "req", signReq.RequestID)

	return nil
}

func (sctUseCase *SignContractRequestUseCase) Accept(request base.IRequest) error {

	if err := request.Validate(sign_contract.SignRequest); err != nil {
		return err
	}
	signReq := request.(*sign_contract.SignContractRequest)

	artist := models.Artist{
		UserID:       signReq.ApplierID,
		Nickname:     signReq.Nickname,
		ContractTerm: cdtime.GetEndOfContract(),
		Activity:     true,
		ManagerID:    signReq.ManagerID,
	}

	ctx := context.Background()
	return sctUseCase.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := sctUseCase.userRepo.UpdateType(ctx, artist.UserID, models.ArtistUser); err != nil {
			sctUseCase.logger.Error("SIGNREQ_UC TRANSACTION Accept", "req", signReq.RequestID, slog.Any("error", err))
			return fmt.Errorf("can't update user with err %w", err)
		}

		if err := sctUseCase.artistRepo.Create(ctx, &artist); err != nil {
			sctUseCase.logger.Error("SIGNREQ_UC TRANSACTION Accept", "req", signReq.RequestID, slog.Any("error", err))
			return fmt.Errorf("can't create artist %s with err %w", artist.Nickname, err)
		}

		signReq.Status = base.ClosedRequest
		if err := sctUseCase.repo.Update(ctx, signReq); err != nil {
			sctUseCase.logger.Error("SIGNREQ_UC TRANSACTION Accept", "req", signReq.RequestID, slog.Any("error", err))
			return fmt.Errorf("can't update reqiest with err %w", err)
		}

		sctUseCase.logger.Debug("SIGNREQ_UC Accept", "req", signReq.RequestID)
		return nil
	})
}

func (sctUseCase *SignContractRequestUseCase) Decline(request base.IRequest) error {

	if err := request.Validate(sign_contract.SignRequest); err != nil {
		return err
	}
	signReq := request.(*sign_contract.SignContractRequest)

	signReq.Status = base.ClosedRequest
	signReq.Description = base.DescrDeclinedRequest

	sctUseCase.logger.Debug("SIGNREQ_UC Decline", "req", signReq.RequestID)

	return sctUseCase.repo.Update(context.Background(), signReq)
}

func (sctUseCase *SignContractRequestUseCase) Get(id uint64) (*sign_contract.SignContractRequest, error) {

	req, err := sctUseCase.repo.Get(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("can't get sign contract request with err %w", err)
	}

	return req, nil
}

func (sctUseCase *SignContractRequestUseCase) sendProceedToManagerMSG(signReq *sign_contract.SignContractRequest) error {

	msg, err := broker_dto.NewSignRequestProducerMsg(signContractBroker.SignRequestProceedToManager, signReq)
	if err != nil {
		return fmt.Errorf("can't apply sign contract request: can't proceed to manager with err %w", err)
	}

	_, _, err = sctUseCase.scBroker.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("can't apply sign contract request: can't proceed to manager with err %w", err)
	}

	return nil
}
