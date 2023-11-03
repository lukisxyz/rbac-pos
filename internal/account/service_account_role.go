package account

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
)

type RoleAccountService interface {
	GetAccount(ctx context.Context, rid ulid.ULID) (AccountRoleList, error)
	GetRoleByAccount(ctx context.Context, uid ulid.ULID) (RoleAccountList, error)
	AssignRole(ctx context.Context, uid, rid ulid.ULID) error
	DeleteRole(ctx context.Context, uid, rid ulid.ULID) error
}

func NewAccountRoleService(
	repo Repo,
	readModel ReadModel,
	accountRoleRepo RepoAccountRole,
	accountRoleReadModel ReadModelAccountRole,
) RoleAccountService {
	return &services{
		repo:                 repo,
		readModel:            readModel,
		accountRoleRepo:      accountRoleRepo,
		accountRoleReadModel: accountRoleReadModel,
	}
}

// AssignRole implements RoleAccountService.
func (s *services) AssignRole(ctx context.Context, rid, uid ulid.ULID) error {
	_, err := s.accountRoleReadModel.Find(ctx, rid, uid)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			return s.accountRoleRepo.AssignRole(ctx, rid, uid)
		}
	}
	return ErrAccountRoleAlreadyAssigned
}

// AssignRole implements RoleAccountService.
func (s *services) DeleteRole(ctx context.Context, rid, uid ulid.ULID) error {
	data, err := s.accountRoleReadModel.Find(ctx, rid, uid)
	if err != nil {
		return err
	}
	return s.accountRoleRepo.RemoveRole(ctx, data.AccountId, data.RoleId)
}

// GetAccount implements RoleAccountService.
func (s *services) GetAccount(ctx context.Context, rid ulid.ULID) (AccountRoleList, error) {
	return s.accountRoleReadModel.FetchByRole(ctx, rid)
}

// GetRoleByAccount implements RoleAccountService.
func (s *services) GetRoleByAccount(ctx context.Context, uid ulid.ULID) (RoleAccountList, error) {
	return s.accountRoleReadModel.FetchByAccount(ctx, uid)
}
