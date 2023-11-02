package role

import (
	"context"
	"errors"
	"pos/permission"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrPermissionAlreadyAssigned = errors.New("role: permission already assigned")
	ErrPermissionNotFound        = errors.New("role: permission not found")
)

type ReadModelRolePermission interface {
	FetchByPermission(ctx context.Context, id ulid.ULID) (PermissionRoleList, error)
	FetchByRole(ctx context.Context, id ulid.ULID) (RolePermissionList, error)
	Find(ctx context.Context, pid, rid ulid.ULID) (*RolePermission, error)
}

type PermissionRoleList struct {
	Roles []Role `json:"data"`
	Count int    `json:"count"`
}

var emptyRole = PermissionRoleList{
	Roles: []Role{},
	Count: 0,
}

// FetchByPermission implements ReadModelRolePermission.
func (r *repo) FetchByPermission(ctx context.Context, id ulid.ULID) (PermissionRoleList, error) {
	var itemCount int
	row := r.db.QueryRow(
		ctx,
		`SELECT COUNT(*) as c FROM role_permissions WHERE permission_id = $1`,
		id,
	)
	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in role permission")
		return emptyRole, err
	}

	if itemCount == 0 {
		return emptyRole, nil
	}
	log.Debug().Int("count", itemCount).Msg("found role permission items")
	items := make([]Role, itemCount)
	rows, err := r.db.Query(
		ctx,
		`
			SELECT
				rp.role_id,
				p.name,
				p.description AS desc
			FROM
				role_permissions rp
			JOIN
				roles p
			ON
				rp.role_id = p.id
			WHERE
				rp.permission_id = $1;
		`,
		id,
	)
	if err != nil {
		return emptyRole, err
	}
	defer rows.Close()

	var count int
	for count = range items {
		var id ulid.ULID
		var name string
		var desc string
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&name,
			&desc,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyRole, err
		}
		items[count] = Role{
			Id:          id,
			Name:        name,
			Description: desc,
			CreatedAt:   time.Time{},
		}
	}
	list := PermissionRoleList{
		Roles: items,
		Count: itemCount,
	}
	return list, nil
}

type RolePermissionList struct {
	Permissions []permission.Permission `json:"data"`
	Count       int                     `json:"count"`
}

var emptyPermission = RolePermissionList{
	Permissions: []permission.Permission{},
	Count:       0,
}

// FetchByRole implements ReadModelRolePermission.
func (r *repo) FetchByRole(ctx context.Context, id ulid.ULID) (RolePermissionList, error) {
	var itemCount int
	row := r.db.QueryRow(
		ctx,
		`SELECT COUNT(*) as c FROM role_permissions WHERE role_id = $1`,
		id,
	)
	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in role permission")
		return emptyPermission, err
	}

	if itemCount == 0 {
		return emptyPermission, nil
	}
	log.Debug().Int("count", itemCount).Msg("found role permission items")
	items := make([]permission.Permission, itemCount)
	rows, err := r.db.Query(
		ctx,
		`
			SELECT
				p.id AS permission_id,
				p.name,
				p.description AS desc,
				p.url
			FROM
				role_permissions rp
			JOIN
				permissions p
			ON
				rp.permission_id = p.id
			WHERE
				rp.role_id = $1;
		`,
		id,
	)
	if err != nil {
		return emptyPermission, err
	}
	defer rows.Close()

	var count int
	for count = range items {
		var id ulid.ULID
		var name string
		var desc string
		var url string
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&name,
			&desc,
			&url,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyPermission, err
		}
		items[count] = permission.Permission{
			Id:          id,
			Name:        name,
			Description: desc,
			Url:         url,
			CreatedAt:   time.Time{},
		}
	}
	list := RolePermissionList{
		Permissions: items,
		Count:       itemCount,
	}
	return list, nil
}

// Find implements ReadModelRolePermission.
func (r *repo) Find(ctx context.Context, pid, rid ulid.ULID) (*RolePermission, error) {
	row := r.db.QueryRow(
		ctx,
		`
			SELECT
				permission_id,
				role_id
			FROM
				role_permissions
			WHERE
				permission_id = $1 AND
				role_id = $2
		`,
		pid,
		rid,
	)
	var data RolePermission
	if err := row.Scan(
		&data.PermissionId,
		&data.RoleId,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Err(err).Msg("can't find any item")
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return &data, nil
}

func NewReadModelRolePermission(db *pgxpool.Pool) ReadModelRolePermission {
	return &repo{db: db}
}

type RepoRolePermission interface {
	AssignPermission(ctx context.Context, roleId, permissionId ulid.ULID) error
	RemovePermission(ctx context.Context, roleId, permissionId ulid.ULID) error
}

// AssignPermission implements RepoRolePermission.
func (r *repo) AssignPermission(ctx context.Context, roleId, permissionId ulid.ULID) error {
	_, err := r.db.Exec(
		ctx,
		`
			INSERT INTO role_permissions (
				role_id,
				permission_id
			) VALUES (
				$1,
				$2
			);
		`,
		roleId,
		permissionId,
	)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == "23505" {
			return ErrPermissionNotFound
		}
		return err
	}
	return nil
}

// RemovePermission implements RepoRolePermission.
func (r *repo) RemovePermission(ctx context.Context, roleId, permissionId ulid.ULID) error {
	_, err := r.db.Exec(
		ctx,
		`
			DELETE FROM role_permissions 
			WHERE permission_id = $1 AND role_id = $2;
		`,
		permissionId,
		roleId,
	)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == "23505" {
			return ErrPermissionNotFound
		}
		return err
	}
	return nil
}

func NewRepoRolePermission(db *pgxpool.Pool) RepoRolePermission {
	return &repo{db: db}
}
