package role

import (
	"context"
	"errors"
	"pos/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRoleNotFound     = errors.New("role: not found")
	ErrRoleAlreadyExist = errors.New("role: url already exists")
)

type repo struct {
	db *pgxpool.Pool
}

// Delete implements Repo.
func (r *repo) Delete(ctx context.Context, data *domain.Role) error {
	_, err := r.db.Exec(
		ctx,
		`
			DELETE FROM roles
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
func (r *repo) Save(ctx context.Context, data *domain.Role) error {
	_, err := r.db.Exec(
		ctx,
		`
			INSERT INTO roles (
				id,
				name,
				description,
				created_at
			) VALUES (
				$1,
				$2,
				$3,
				$4
			) ON CONFLICT (id) DO UPDATE
			SET name = excluded.name,
				description = excluded.description;
		`,
		data.Id,
		data.Name,
		data.Description,
		data.CreatedAt,
	)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == "23505" {
			return ErrRoleAlreadyExist
		}
		return err
	}
	return nil
}

type Repo interface {
	Save(ctx context.Context, data *domain.Role) error
	Delete(ctx context.Context, data *domain.Role) error
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{db: db}
}
