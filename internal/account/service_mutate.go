package account

import (
	"context"
	"pos/domain"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

// CreateAccount implements MutationData.
func (s *services) CreateAccount(ctx context.Context, email, pwd string) (*domain.Account, error) {
	newData, err := domain.NewAccount(email, pwd)
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
func (s *services) EditAccount(ctx context.Context, id ulid.ULID, pwd string) (*domain.Account, error) {
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
	CreateAccount(ctx context.Context, email, pwd string) (*domain.Account, error)
	EditAccount(ctx context.Context, id ulid.ULID, pwd string) (*domain.Account, error)
	DeleteAccount(ctx context.Context, id ulid.ULID) error
}

func NewMutationData(
	repo Repo,
	readModel ReadModel,
) MutationData {
	return &services{repo: repo, readModel: readModel}
}
