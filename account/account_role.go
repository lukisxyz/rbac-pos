package account

import (
	"context"
	"errors"
	"fmt"
	"pos/role"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrRoleNotFound               = errors.New("account role: not found")
	ErrAccountRoleAlreadyAssigned = errors.New("account role: already assigned")
)

type ReadModelAccountRole interface {
	FetchByAccount(ctx context.Context, id ulid.ULID) (RoleAccountList, error)
	FetchByRole(ctx context.Context, id ulid.ULID) (AccountRoleList, error)
	Find(ctx context.Context, rid, uid ulid.ULID) (*AccountRole, error)
}

type RoleAccountList struct {
	Roles []role.Role `json:"data"`
	Count int         `json:"count"`
}

var emptyRole = RoleAccountList{
	Roles: []role.Role{},
	Count: 0,
}

// FetchByAccount implements ReadModelAccountRole.
func (r *repo) FetchByAccount(ctx context.Context, id ulid.ULID) (RoleAccountList, error) {
	var itemCount int
	row := r.db.QueryRow(
		ctx,
		`SELECT COUNT(*) as c FROM account_roles WHERE account_id = $1`,
		id,
	)
	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in role")
		return emptyRole, err
	}

	if itemCount == 0 {
		return emptyRole, nil
	}
	log.Debug().Int("count", itemCount).Msg("found role items")
	items := make([]role.Role, itemCount)
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			p.id,
			p.name AS role_name,
			p.description AS role_description
		FROM
			account_roles rp
		JOIN
			roles p
		ON
			rp.role_id = p.id
		WHERE
			rp.account_id = $1;	
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
		fmt.Println(name)
		items[count] = role.Role{
			Id:          id,
			Name:        name,
			Description: desc,
			CreatedAt:   time.Time{},
		}
	}
	list := RoleAccountList{
		Roles: items,
		Count: itemCount,
	}
	return list, nil
}

type AccountRoleList struct {
	Accounts []Account `json:"data"`
	Count    int       `json:"count"`
}

var emptyAccount = AccountRoleList{
	Accounts: []Account{},
	Count:    0,
}

// FetchByRole implements ReadModelAccountRole.
func (r *repo) FetchByRole(ctx context.Context, id ulid.ULID) (AccountRoleList, error) {
	var itemCount int
	row := r.db.QueryRow(
		ctx,
		`SELECT COUNT(*) as c FROM account_roles WHERE role_id = $1`,
		id,
	)
	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in role")
		return emptyAccount, err
	}

	if itemCount == 0 {
		return emptyAccount, nil
	}
	log.Debug().Int("count", itemCount).Msg("found role items")
	items := make([]Account, itemCount)
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			p.id AS account_id,
			p.email
		FROM
			account_roles rp
		JOIN
			accounts p
		ON
			rp.account_id = p.id
		WHERE
			rp.role_id = $1;
		`,
		id,
	)
	if err != nil {
		return emptyAccount, err
	}
	defer rows.Close()

	var count int
	for count = range items {
		var id ulid.ULID
		var email string
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&email,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyAccount, err
		}
		items[count] = Account{
			Id:    id,
			Email: email,
		}
	}
	list := AccountRoleList{
		Accounts: items,
		Count:    itemCount,
	}
	return list, nil
}

// Find implements ReadModelAccountRole.
func (r *repo) Find(ctx context.Context, rid, uid ulid.ULID) (*AccountRole, error) {
	row := r.db.QueryRow(
		ctx,
		`
			SELECT
				role_id,
				account_id
			FROM
				account_roles
			WHERE
				role_id = $1 AND
				account_id = $2
		`,
		rid,
		uid,
	)
	var data AccountRole
	if err := row.Scan(
		&data.RoleId,
		&data.AccountId,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Err(err).Msg("can't find any item")
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &data, nil
}

func NewReadModelAccountRole(db *pgxpool.Pool) ReadModelAccountRole {
	return &repo{db: db}
}

type RepoAccountRole interface {
	AssignRole(ctx context.Context, accountId, roleId ulid.ULID) error
	RemoveRole(ctx context.Context, accountId, roleId ulid.ULID) error
}

// AssignRole implements RepoAccountRole.
func (r *repo) AssignRole(ctx context.Context, roleId, accountId ulid.ULID) error {
	_, err := r.db.Exec(
		ctx,
		`
			INSERT INTO account_roles (
				account_id,
				role_id
			) VALUES (
				$1,
				$2
			);
		`,
		accountId,
		roleId,
	)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == "23505" {
			return ErrRoleNotFound
		}
		return err
	}
	return nil
}

// RemoveRole implements RepoAccountRole.
func (r *repo) RemoveRole(ctx context.Context, accountId, roleId ulid.ULID) error {
	_, err := r.db.Exec(
		ctx,
		`
			DELETE FROM account_roles 
			WHERE account_id = $1 AND role_id = $2;
		`,
		accountId,
		roleId,
	)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == "23505" {
			return ErrRoleNotFound
		}
		return err
	}
	return nil
}

func NewRepoAccountRole(db *pgxpool.Pool) RepoAccountRole {
	return &repo{db: db}
}
