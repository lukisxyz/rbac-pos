package account

import (
	"context"
	"errors"
	"pos/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAccountNotFound     = errors.New("account: not found")
	ErrAccountAlreadyExist = errors.New("account: email already exists")
)

type repo struct {
	db *pgxpool.Pool
}

// Delete implements Repo.
func (r *repo) Delete(ctx context.Context, data *domain.Account) error {
	_, err := r.db.Exec(
		ctx,
		`
			DELETE FROM accounts
			WHERE id = $1
		`,
		data.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

// Save implements Repo.
func (r *repo) Save(ctx context.Context, data *domain.Account) error {
	_, err := r.db.Exec(
		ctx,
		`
			INSERT INTO accounts (
				id,
				email,
				password,
				created_at
			) VALUES (
				$1,
				$2,
				$3,
				$4
			) ON CONFLICT (id) DO UPDATE
			SET
				password = excluded.password;
		`,
		data.Id,
		data.Email,
		data.Password,
		data.CreatedAt,
	)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == "23505" {
			return ErrAccountAlreadyExist
		}
		return err
	}
	return nil
}

type Repo interface {
	Save(ctx context.Context, data *domain.Account) error
	Delete(ctx context.Context, data *domain.Account) error
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{db: db}
}
