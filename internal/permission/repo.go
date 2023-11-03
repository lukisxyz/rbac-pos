package permission

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
	ErrPermissionNotFound     = errors.New("permission: not found")
	ErrPermissionAlreadyExist = errors.New("permission: url already exists")
)

type repo struct {
	db *pgxpool.Pool
}

type PermissionList struct {
	Permissions []Permission `json:"data"`
	Count       int          `json:"count"`
}

var emptyList = PermissionList{
	Permissions: []Permission{},
	Count:       0,
}

// Fetch implements ReadModel.
func (r *repo) Fetch(ctx context.Context) (PermissionList, error) {
	var itemCount int

	row := r.db.QueryRow(
		ctx,
		`SELECT COUNT(id) as c FROM permissions`,
	)
	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in permission")
		return emptyList, err
	}

	if itemCount == 0 {
		return emptyList, nil
	}
	log.Debug().Int("count", itemCount).Msg("found permission items")
	items := make([]Permission, itemCount)
	rows, err := r.db.Query(
		ctx,
		`
			SELECT
				id,
				name,
				description,
				url,
				created_at
			FROM
				permissions
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
		var url string
		var createdAt time.Time
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&name,
			&desc,
			&url,
			&createdAt,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyList, err
		}
		items[count] = Permission{
			Id:          id,
			Name:        name,
			Description: desc,
			Url:         url,
			CreatedAt:   createdAt,
		}
	}
	list := PermissionList{
		Permissions: items,
		Count:       itemCount,
	}
	return list, nil
}

// FindById implements ReadModel.
func (r *repo) FindById(ctx context.Context, id ulid.ULID) (*Permission, error) {
	row := r.db.QueryRow(
		ctx,
		`
			SELECT
				id,
				name,
				description,
				url,
				created_at
			FROM
				permissions
			WHERE
				id = $1
		`,
		id,
	)
	var data Permission
	if err := row.Scan(
		&data.Id,
		&data.Name,
		&data.Description,
		&data.Url,
		&data.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Err(err).Msg("can't find any item")
			return nil, ErrPermissionNotFound
		}
	}
	return &data, nil
}

// FindByUrl implements ReadModel.
func (*repo) FindByUrl(ctx context.Context, url string) (*Permission, error) {
	panic("unimplemented")
}

// Delete implements Repo.
func (r *repo) Delete(ctx context.Context, data *Permission) error {
	_, err := r.db.Exec(
		ctx,
		`
			DELETE FROM permissions
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
func (r *repo) Save(ctx context.Context, data *Permission) error {
	_, err := r.db.Exec(
		ctx,
		`
			INSERT INTO permissions (
				id,
				name,
				description,
				url,
				created_at
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5
			) ON CONFLICT (id) DO UPDATE
			SET name = excluded.name,
				description = excluded.description,
				url = excluded.url;
		`,
		data.Id,
		data.Name,
		data.Description,
		data.Url,
		data.CreatedAt,
	)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == "23505" {
			return ErrPermissionAlreadyExist
		}
		return err
	}
	return nil
}

type Repo interface {
	Save(ctx context.Context, data *Permission) error
	Delete(ctx context.Context, data *Permission) error
}

type ReadModel interface {
	Fetch(ctx context.Context) (PermissionList, error)
	FindById(ctx context.Context, id ulid.ULID) (*Permission, error)
	FindByUrl(ctx context.Context, url string) (*Permission, error)
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{db: db}
}

func NewReadModel(db *pgxpool.Pool) ReadModel {
	return &repo{db: db}
}
