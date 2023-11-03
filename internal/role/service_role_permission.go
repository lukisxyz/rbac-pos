package role

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
)

type RolePermissionService interface {
	GetPermission(ctx context.Context, rid ulid.ULID) (RolePermissionList, error)
	GetRoleByPermission(ctx context.Context, pid ulid.ULID) (PermissionRoleList, error)
	AssignPermisson(ctx context.Context, uid, rid, pid ulid.ULID) error
	DeletePermission(ctx context.Context, uid, rid, pid ulid.ULID) error
}

func NewRolePermissionService(
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
func (s *services) DeletePermission(ctx context.Context, uid, rid, pid ulid.ULID) error {
	data, err := s.rolePermissionReadModel.Find(ctx, pid, rid)
	if err != nil {
		return err
	}
	err = s.rolePermissionRepo.RemovePermission(ctx, data.RoleId, data.PermissionId)
	if err != nil {
		return err
	}
	return s.rolePermissionRepo.RevokeSession(ctx, uid)
}

// AssignPermisson implements RolePermissionService.
func (s *services) AssignPermisson(ctx context.Context, uid, rid, pid ulid.ULID) error {
	_, err := s.rolePermissionReadModel.Find(ctx, pid, rid)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			err := s.rolePermissionRepo.AssignPermission(ctx, rid, pid)
			if err != nil {
				return err
			}
			return s.rolePermissionRepo.RevokeSession(ctx, uid)
		}
	}
	return ErrPermissionAlreadyAssigned
}

// GetPermission implements RolePermissionService.
func (s *services) GetPermission(ctx context.Context, rid ulid.ULID) (RolePermissionList, error) {
	return s.rolePermissionReadModel.FetchByRole(ctx, rid)
}

// GetRoleByPermission implements RolePermissionService.
func (s *services) GetRoleByPermission(ctx context.Context, pid ulid.ULID) (PermissionRoleList, error) {
	return s.rolePermissionReadModel.FetchByPermission(ctx, pid)
}
