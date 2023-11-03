package role

import (
	"context"
	"pos/domain"

	"github.com/oklog/ulid/v2"
)

// CreateRole implements MutationData.
func (s *services) CreateRole(ctx context.Context, name, desc string) (*domain.Role, error) {
	newData := domain.NewRole(name, desc)
	if err := s.repo.Save(ctx, &newData); err != nil {
		return nil, err
	}
	return &newData, nil
}

// DeleteRole implements MutationData.
func (s *services) DeleteRole(ctx context.Context, id ulid.ULID) error {
	currentData, err := s.readModel.FindById(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, currentData)
}

// EditRole implements MutationData.
func (s *services) EditRole(ctx context.Context, id ulid.ULID, name, desc string) (*domain.Role, error) {
	currentData, err := s.readModel.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	currentData.Name = name
	currentData.Description = desc
	if err := s.repo.Save(ctx, currentData); err != nil {
		return nil, err
	}
	return currentData, nil
}

type MutationData interface {
	CreateRole(ctx context.Context, name, desc string) (*domain.Role, error)
	EditRole(ctx context.Context, id ulid.ULID, name, desc string) (*domain.Role, error)
	DeleteRole(ctx context.Context, id ulid.ULID) error
}

func NewMutationData(
	repo Repo,
	readModel ReadModel,
) MutationData {
	return &services{repo: repo, readModel: readModel}
}
