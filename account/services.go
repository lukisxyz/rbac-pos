package account

import (
	"context"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type services struct {
	repo      Repo
	readModel ReadModel
}

// GetAll implements ReadData.
func (s *services) GetAll(ctx context.Context) (AccountList, error) {
	return s.readModel.Fetch(ctx)
}

// GetOneById implements ReadData.
func (s *services) GetOneById(ctx context.Context, id ulid.ULID) (*Account, error) {
	return s.readModel.FindById(ctx, id)
}

// CreateAccount implements MutationData.
func (s *services) CreateAccount(ctx context.Context, email, pwd string) (*Account, error) {
	newData, err := newAccount(email, pwd)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, &newData); err != nil {
		return nil, err
	}
	return &newData, nil
}

// DeleteAccount implements MutationData.
func (s *services) DeleteAccount(ctx context.Context, id ulid.ULID) error {
	currentData, err := s.readModel.FindById(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, currentData)
}

// EditAccount implements MutationData.
func (s *services) EditAccount(ctx context.Context, id ulid.ULID, pwd string) (*Account, error) {
	currentData, err := s.readModel.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	currentData.Password = string(bytes)
	if err := s.repo.Save(ctx, currentData); err != nil {
		return nil, err
	}
	return currentData, nil
}

type MutationData interface {
	CreateAccount(ctx context.Context, email, pwd string) (*Account, error)
	EditAccount(ctx context.Context, id ulid.ULID, pwd string) (*Account, error)
	DeleteAccount(ctx context.Context, id ulid.ULID) error
}

func NewMutationData(
	repo Repo,
	readModel ReadModel,
) MutationData {
	return &services{repo: repo, readModel: readModel}
}

type ReadData interface {
	GetAll(ctx context.Context) (AccountList, error)
	GetOneById(ctx context.Context, id ulid.ULID) (*Account, error)
}

func NewReadData(
	repo Repo,
	readModel ReadModel,
) ReadData {
	return &services{repo: repo, readModel: readModel}
}
