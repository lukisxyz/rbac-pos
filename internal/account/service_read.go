package account

import (
	"context"
	"pos/domain"

	"github.com/oklog/ulid/v2"
)

type services struct {
	repo                 Repo
	readModel            ReadModel
	accountRoleRepo      RepoAccountRole
	accountRoleReadModel ReadModelAccountRole
}

// GetAll implements ReadData.
func (s *services) GetAll(ctx context.Context) (AccountList, error) {
	return s.readModel.Fetch(ctx)
}

// GetOneById implements ReadData.
func (s *services) GetOneById(ctx context.Context, id ulid.ULID) (*domain.Account, error) {
	return s.readModel.FindById(ctx, id)
}

type ReadData interface {
	GetAll(ctx context.Context) (AccountList, error)
	GetOneById(ctx context.Context, id ulid.ULID) (*domain.Account, error)
}

func NewReadData(
	repo Repo,
	readModel ReadModel,
) ReadData {
	return &services{repo: repo, readModel: readModel}
}
