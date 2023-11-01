package role

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrRoleNotFound     = errors.New("role: not found")
	ErrRoleAlreadyExist = errors.New("role: url already exists")
)

type repo struct {
	db *pgxpool.Pool
}

type RoleList struct {
	Roles []Role `json:"data"`
	Count int    `json:"count"`
}

var emptyList = RoleList{
	Roles: []Role{},
	Count: 0,
}

// Fetch implements ReadModel.
func (r *repo) Fetch(ctx context.Context) (RoleList, error) {
	var itemCount int

	row := r.db.QueryRow(
		ctx,
		`SELECT COUNT(id) as c FROM roles`,
	)
	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in Role")
		return emptyList, err
	}

	if itemCount == 0 {
		return emptyList, nil
	}
	log.Debug().Int("count", itemCount).Msg("found Role items")
	items := make([]Role, itemCount)
	rows, err := r.db.Query(
		ctx,
		`
			SELECT
				id,
				name,
				description,
				created_at
			FROM
				roles
			ORDER BY
				id
		`,
	)
	if err != nil {
		return emptyList, err
	}
	defer rows.Close()

	var count int
	for count = range items {
		var id ulid.ULID
		var name string
		var desc string
		var createdAt time.Time
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&name,
			&desc,
			&createdAt,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyList, err
		}
		items[count] = Role{
			Id:          id,
			Name:        name,
			Description: desc,
			CreatedAt:   createdAt,
		}
	}
	list := RoleList{
		Roles: items,
		Count: count,
	}
	return list, nil
}

// FindById implements ReadModel.
func (r *repo) FindById(ctx context.Context, id ulid.ULID) (*Role, error) {
	row := r.db.QueryRow(
		ctx,
		`
			SELECT
				id,
				name,
				description,
				created_at
			FROM
				roles
			WHERE
				id = $1
		`,
		id,
	)
	var data Role
	if err := row.Scan(
		&data.Id,
		&data.Name,
		&data.Description,
		&data.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Err(err).Msg("can't find any item")
			return nil, ErrRoleNotFound
		}
	}
	return &data, nil
}

// FindByUrl implements ReadModel.
func (*repo) FindByUrl(ctx context.Context, url string) (*Role, error) {
	panic("unimplemented")
}

// Delete implements Repo.
func (r *repo) Delete(ctx context.Context, data *Role) error {
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
func (r *repo) Save(ctx context.Context, data *Role) error {
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
	Save(ctx context.Context, data *Role) error
	Delete(ctx context.Context, data *Role) error
}

type ReadModel interface {
	Fetch(ctx context.Context) (RoleList, error)
	FindById(ctx context.Context, id ulid.ULID) (*Role, error)
	FindByUrl(ctx context.Context, url string) (*Role, error)
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{db: db}
}

func NewReadModel(db *pgxpool.Pool) ReadModel {
	return &repo{db: db}
}
