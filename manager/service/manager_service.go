package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rauzh/cd-core/models"

	"github.com/rauzh/cd-core/repo"
)

type IManagerService interface {
	Create(*models.Manager) error
	Get(uint64) (*models.Manager, error)
	GetByUserID(uint64) (*models.Manager, error)
	GetRandomManagerID() (uint64, error)
	GetForAdmin() ([]models.Manager, error)
}

type ManagerService struct {
	repo   repo.ManagerRepo
	logger *slog.Logger
}

func NewManagerService(r repo.ManagerRepo, logger *slog.Logger) IManagerService {
	return &ManagerService{repo: r, logger: logger}
}

func (mngSvc *ManagerService) Create(manager *models.Manager) error {
	if err := mngSvc.repo.Create(context.Background(), manager); err != nil {
		mngSvc.logger.Error("MANAGER_SERVICE Create", "error", err.Error())
		return fmt.Errorf("can't create manager with err %w", err)
	}
	mngSvc.logger.Info("MANAGER_SERVICE Create", "manager_userID", manager.UserID)
	return nil
}

func (mngSvc *ManagerService) Get(id uint64) (*models.Manager, error) {
	manager, err := mngSvc.repo.Get(context.Background(), id)

	if err != nil {
		mngSvc.logger.Error("MANAGER_SERVICE Get", "error", err.Error())
		return nil, fmt.Errorf("can't get manager with err %w", err)
	}
	mngSvc.logger.Debug("MANAGER_SERVICE Get", "manager_userID", manager.UserID)
	return manager, nil
}

func (mngSvc *ManagerService) GetForAdmin() ([]models.Manager, error) {
	managers, err := mngSvc.repo.GetForAdmin(context.Background())

	if err != nil {
		mngSvc.logger.Error("MANAGER SERVICE GetForAdmin", "error", err.Error())
		return nil, fmt.Errorf("can't get managers with err %w", err)
	}
	mngSvc.logger.Debug("MANAGER_SERVICE GetForAdmin", "managers_len", len(managers))
	return managers, nil
}

func (mngSvc *ManagerService) GetByUserID(id uint64) (*models.Manager, error) {
	manager, err := mngSvc.repo.GetByUserID(context.Background(), id)

	if err != nil {
		mngSvc.logger.Error("MANAGER_SERVICE GetByUserID", "error", err.Error())
		return nil, fmt.Errorf("can't get manager with err %w", err)
	}
	mngSvc.logger.Debug("MANAGER_SERVICE Get", "manager_userID", manager.UserID)
	return manager, nil
}

func (mngSvc *ManagerService) GetRandomManagerID() (uint64, error) {
	id, err := mngSvc.repo.GetRandManagerID(context.Background())
	if err != nil {
		mngSvc.logger.Error("MANAGER_SERVICE GetRandomManagerID", "error", err.Error())
		return 0, fmt.Errorf("can't get manager with err %w", err)
	}
	mngSvc.logger.Debug("MANAGER_SERVICE Get", "rand_managerID", id)
	return id, nil
}
