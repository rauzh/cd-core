package transactor

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/rauzh/cd-core/transactor"
)

type ATtrm struct {
	trm *manager.Manager
}

func NewATtrm(trm *manager.Manager) transactor.Transactor {
	return &ATtrm{trm: trm}
}

func (transactor *ATtrm) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return transactor.trm.Do(ctx, fn)
}
