package role

import (
	"context"
	"pos/domain"

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
func (s *services) GetOneById(ctx context.Context, id ulid.ULID) (*domain.ReadRoleResponse, error) {
	role, err := s.readModel.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	permissionList, err := s.rolePermissionReadModel.FetchByRole(ctx, id)
	if err != nil {
		return nil, err
	}
	data := domain.ReadRoleResponse{
		Role:             *role,
		TotalPermissions: permissionList.Count,
		Permissions:      permissionList.Permissions,
	}
	return &data, nil
}

type ReadData interface {
	GetAll(ctx context.Context) (RoleList, error)
	GetOneById(ctx context.Context, id ulid.ULID) (*domain.ReadRoleResponse, error)
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
