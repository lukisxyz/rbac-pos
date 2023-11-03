package role

import (
	"context"
	"errors"
	"pos/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type RoleList struct {
	Roles []domain.Role `json:"data"`
	Count int           `json:"count"`
}

var emptyList = RoleList{
	Roles: []domain.Role{},
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
	items := make([]domain.Role, itemCount)
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
		items[count] = domain.Role{
			Id:          id,
			Name:        name,
			Description: desc,
			CreatedAt:   createdAt,
		}
	}
	list := RoleList{
		Roles: items,
		Count: itemCount,
	}
	return list, nil
}

// FindById implements ReadModel.
func (r *repo) FindById(ctx context.Context, id ulid.ULID) (*domain.Role, error) {
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
	var data domain.Role
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

type ReadModel interface {
	Fetch(ctx context.Context) (RoleList, error)
	FindById(ctx context.Context, id ulid.ULID) (*domain.Role, error)
}

func NewReadModel(db *pgxpool.Pool) ReadModel {
	return &repo{db: db}
}
