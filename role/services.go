package role

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
)

type services struct {
	repo                    Repo
	readModel               ReadModel
	rolePermissionRepo      RepoRolePermission
	rolePermissionReadModel ReadModelRolePermission
}

// GetAll implements ReadData.
func (s *services) GetAll(ctx context.Context) (RoleList, error) {
	return s.readModel.Fetch(ctx)
}

// GetOneById implements ReadData.
func (s *services) GetOneById(ctx context.Context, id ulid.ULID) (*Role, error) {
	return s.readModel.FindById(ctx, id)
}

// CreateRole implements MutationData.
func (s *services) CreateRole(ctx context.Context, name, desc string) (*Role, error) {
	newData := newRole(name, desc)
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
func (s *services) EditRole(ctx context.Context, id ulid.ULID, name, desc string) (*Role, error) {
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
	CreateRole(ctx context.Context, name, desc string) (*Role, error)
	EditRole(ctx context.Context, id ulid.ULID, name, desc string) (*Role, error)
	DeleteRole(ctx context.Context, id ulid.ULID) error
}

func NewMutationData(
	repo Repo,
	readModel ReadModel,
) MutationData {
	return &services{repo: repo, readModel: readModel}
}

type ReadData interface {
	GetAll(ctx context.Context) (RoleList, error)
	GetOneById(ctx context.Context, id ulid.ULID) (*Role, error)
}

func NewReadData(
	repo Repo,
	readModel ReadModel,
	rolePermissionRepo RepoRolePermission,
	rolePermissionReadModel ReadModelRolePermission,
) ReadData {
	return &services{
		repo:                    repo,
		readModel:               readModel,
		rolePermissionRepo:      rolePermissionRepo,
		rolePermissionReadModel: rolePermissionReadModel,
	}
}

type RolePermissionService interface {
	GetPermission(ctx context.Context, rid ulid.ULID) (RolePermissionList, error)
	AssignPermisson(ctx context.Context, rid, pid ulid.ULID) error
}

func NewRelePermissionService(
	repo Repo,
	readModel ReadModel,
	rolePermissionRepo RepoRolePermission,
	rolePermissionReadModel ReadModelRolePermission,
) RolePermissionService {
	return &services{
		repo:                    repo,
		readModel:               readModel,
		rolePermissionRepo:      rolePermissionRepo,
		rolePermissionReadModel: rolePermissionReadModel,
	}
}

// AssignPermisson implements RolePermissionService.
func (s *services) AssignPermisson(ctx context.Context, rid, pid ulid.ULID) error {
	_, err := s.rolePermissionReadModel.Find(ctx, pid, rid)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			return s.rolePermissionRepo.AssignPermission(ctx, rid, pid)
		}
	}
	return ErrPermissionAlreadyAssigned
}

// GetPermission implements RolePermissionService.
func (s *services) GetPermission(ctx context.Context, rid ulid.ULID) (RolePermissionList, error) {
	return s.rolePermissionReadModel.FetchByRole(ctx, rid)
}
