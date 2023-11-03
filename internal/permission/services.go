package permission

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type services struct {
	repo      Repo
	readModel ReadModel
}

// GetAll implements ReadData.
func (s *services) GetAll(ctx context.Context) (PermissionList, error) {
	return s.readModel.Fetch(ctx)
}

// GetOneById implements ReadData.
func (s *services) GetOneById(ctx context.Context, id ulid.ULID) (*Permission, error) {
	return s.readModel.FindById(ctx, id)
}

// GetOneByUrl implements ReadData.
func (*services) GetOneByUrl(ctx context.Context, url string) (*Permission, error) {
	panic("unimplemented")
}

// CreatePermission implements MutationData.
func (s *services) CreatePermission(ctx context.Context, name, desc, url string) (*Permission, error) {
	newData := newPermission(name, desc, url)
	if err := s.repo.Save(ctx, &newData); err != nil {
		return nil, err
	}
	return &newData, nil
}

// DeletePermission implements MutationData.
func (s *services) DeletePermission(ctx context.Context, id ulid.ULID) error {
	currentData, err := s.readModel.FindById(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, currentData)
}

// EditPermission implements MutationData.
func (s *services) EditPermission(ctx context.Context, id ulid.ULID, name, desc, url string) (*Permission, error) {
	currentData, err := s.readModel.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	currentData.Name = name
	currentData.Description = desc
	currentData.Url = url
	if err := s.repo.Save(ctx, currentData); err != nil {
		return nil, err
	}
	return currentData, nil
}

type MutationData interface {
	CreatePermission(ctx context.Context, name, desc, url string) (*Permission, error)
	EditPermission(ctx context.Context, id ulid.ULID, name, desc, url string) (*Permission, error)
	DeletePermission(ctx context.Context, id ulid.ULID) error
}

func NewMutationData(
	repo Repo,
	readModel ReadModel,
) MutationData {
	return &services{repo: repo, readModel: readModel}
}

type ReadData interface {
	GetAll(ctx context.Context) (PermissionList, error)
	GetOneById(ctx context.Context, id ulid.ULID) (*Permission, error)
	GetOneByUrl(ctx context.Context, url string) (*Permission, error)
}

func NewReadData(
	repo Repo,
	readModel ReadModel,
) ReadData {
	return &services{repo: repo, readModel: readModel}
}
