package account

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
	ErrAccountNotFound     = errors.New("account: not found")
	ErrAccountAlreadyExist = errors.New("account: email already exists")
)

type repo struct {
	db *pgxpool.Pool
}

type AccountList struct {
	Accounts []Account `json:"data"`
	Count    int       `json:"count"`
}

var emptyList = AccountList{
	Accounts: []Account{},
	Count:    0,
}

// Fetch implements ReadModel.
func (r *repo) Fetch(ctx context.Context) (AccountList, error) {
	var itemCount int

	row := r.db.QueryRow(
		ctx,
		`SELECT COUNT(id) as c FROM accounts`,
	)
	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in Account")
		return emptyList, err
	}

	if itemCount == 0 {
		return emptyList, nil
	}
	log.Debug().Int("count", itemCount).Msg("found Account items")
	items := make([]Account, itemCount)
	rows, err := r.db.Query(
		ctx,
		`
			SELECT
				id,
				email,
				password,
				created_at
			FROM
				accounts
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
		var email string
		var password string
		var createdAt time.Time
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&email,
			&password,
			&createdAt,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyList, err
		}
		items[count] = Account{
			Id:        id,
			Password:  password,
			Email:     email,
			CreatedAt: createdAt,
		}
	}
	list := AccountList{
		Accounts: items,
		Count:    itemCount,
	}
	return list, nil
}

// FindById implements ReadModel.
func (r *repo) FindById(ctx context.Context, id ulid.ULID) (*Account, error) {
	row := r.db.QueryRow(
		ctx,
		`
			SELECT
				id,
				email,
				password,
				created_at
			FROM
				accounts
			WHERE
				id = $1
		`,
		id,
	)
	var data Account
	if err := row.Scan(
		&data.Id,
		&data.Email,
		&data.Password,
		&data.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Err(err).Msg("can't find any item")
			return nil, ErrAccountNotFound
		}
	}
	return &data, nil
}

// FindByEmail implements ReadModel.
func (r *repo) FindByEmail(ctx context.Context, email string) (*Account, error) {
	row := r.db.QueryRow(
		ctx,
		`
			SELECT
				id,
				email,
				password,
				created_at
			FROM
				accounts
			WHERE
				email = $1
		`,
		email,
	)
	var data Account
	if err := row.Scan(
		&data.Id,
		&data.Email,
		&data.Password,
		&data.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Err(err).Msg("can't find any item")
			return nil, ErrAccountNotFound
		}
	}
	return &data, nil
}

// Delete implements Repo.
func (r *repo) Delete(ctx context.Context, data *Account) error {
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
func (r *repo) Save(ctx context.Context, data *Account) error {
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
	Save(ctx context.Context, data *Account) error
	Delete(ctx context.Context, data *Account) error
}

type ReadModel interface {
	Fetch(ctx context.Context) (AccountList, error)
	FindById(ctx context.Context, id ulid.ULID) (*Account, error)
	FindByEmail(ctx context.Context, email string) (*Account, error)
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{db: db}
}

func NewReadModel(db *pgxpool.Pool) ReadModel {
	return &repo{db: db}
}
