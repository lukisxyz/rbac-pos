package oauth

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

var (
	ErrRefreshTokenNotFound     = errors.New("refresh token: not found")
	ErrRefreshTokenAlreadyExist = errors.New("refresh token: url already exists")
)

type repo struct {
	db *pgxpool.Pool
}

type RefreshTokenList struct {
	RefreshTokens []domain.RefreshToken `json:"data"`
	Count         int                   `json:"count"`
}

var emptyList = RefreshTokenList{
	RefreshTokens: []domain.RefreshToken{},
	Count:         0,
}

// Fetch implements ReadModel.
func (r *repo) Fetch(ctx context.Context) (RefreshTokenList, error) {
	var itemCount int

	row := r.db.QueryRow(
		ctx,
		`SELECT COUNT(id) as c FROM refresh_tokens`,
	)
	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in RefreshToken")
		return emptyList, err
	}

	if itemCount == 0 {
		return emptyList, nil
	}
	log.Debug().Int("count", itemCount).Msg("found RefreshToken items")
	items := make([]domain.RefreshToken, itemCount)
	rows, err := r.db.Query(
		ctx,
		`
			SELECT
				id,
				token_value,
				account_id,
				created_at,
				expires_at,
				revoked
			FROM
				refresh_tokens
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
		var tokenValue string
		var userId ulid.ULID
		var createdAt time.Time
		var expiresAt time.Time
		var revoked bool
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&tokenValue,
			&userId,
			&createdAt,
			&expiresAt,
			&revoked,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyList, err
		}
		items[count] = domain.RefreshToken{
			ID:         id,
			TokenValue: tokenValue,
			UserID:     userId,
			ExpiresAt:  expiresAt,
			Revoked:    revoked,
			CreatedAt:  createdAt,
		}
	}
	list := RefreshTokenList{
		RefreshTokens: items,
		Count:         itemCount,
	}
	return list, nil
}

// FindById implements ReadModel.
func (r *repo) FindById(ctx context.Context, id ulid.ULID) (*domain.RefreshToken, error) {
	query := `
		SELECT
			id,
			token_value,
			account_id,
			created_at,
			expires_at,
			revoked
		FROM
			refresh_tokens
		WHERE
			id = $1
			AND expires_at > NOW()  
			AND revoked = false;   
	`
	row := r.db.QueryRow(
		ctx,
		query,
		id,
	)
	var item domain.RefreshToken
	if err := row.Scan(
		&item.ID,
		&item.TokenValue,
		&item.UserID,
		&item.CreatedAt,
		&item.ExpiresAt,
		&item.Revoked,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return &domain.RefreshToken{}, ErrRefreshTokenNotFound
		}
		return &domain.RefreshToken{}, err
	}
	return &item, nil
}

// FindById implements ReadModel.
func (r *repo) FindByUserID(ctx context.Context, id ulid.ULID) (*domain.RefreshToken, error) {
	query := `
		SELECT
			id,
			token_value,
			account_id,
			created_at,
			expires_at,
			revoked
		FROM
			refresh_tokens
		WHERE
			account_id = $1
			AND expires_at > NOW()  
			AND revoked = false;   
	`
	row := r.db.QueryRow(
		ctx,
		query,
		id,
	)
	var item domain.RefreshToken
	if err := row.Scan(
		&item.ID,
		&item.TokenValue,
		&item.UserID,
		&item.CreatedAt,
		&item.ExpiresAt,
		&item.Revoked,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return &domain.RefreshToken{}, ErrRefreshTokenNotFound
		}
		return &domain.RefreshToken{}, err
	}
	return &item, nil
}

// FindById implements ReadModel.
func (r *repo) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	query := `
		SELECT
			id,
			token_value,
			account_id,
			created_at,
			expires_at,
			revoked
		FROM
			refresh_tokens
		WHERE
			token_value = $1
			AND expires_at > NOW()  
			AND revoked = false;   
	`
	row := r.db.QueryRow(
		ctx,
		query,
		token,
	)
	var item domain.RefreshToken
	if err := row.Scan(
		&item.ID,
		&item.TokenValue,
		&item.UserID,
		&item.CreatedAt,
		&item.ExpiresAt,
		&item.Revoked,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return &domain.RefreshToken{}, ErrRefreshTokenNotFound
		}
		return &domain.RefreshToken{}, err
	}
	return &item, nil
}

// Save implements Repo.
func (r *repo) Save(ctx context.Context, data *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens
			(id, token_value, account_id, created_at, expires_at, revoked)
		VALUES
			($1, $2, $3, $4, $5, $6);
	`

	if _, err := r.db.Exec(
		ctx,
		query,
		data.ID,
		data.TokenValue,
		data.UserID,
		data.CreatedAt,
		data.ExpiresAt,
		data.Revoked,
	); err != nil {
		return err
	}

	return nil
}

// Save implements Repo.
func (r *repo) Revoke(ctx context.Context, data *domain.RefreshToken) error {
	query := `
		UPDATE refresh_tokens
		SET
			revoked = TRUE
		WHERE id = $1
	`

	if _, err := r.db.Exec(
		ctx,
		query,
		data.ID,
	); err != nil {
		return err
	}

	return nil
}

type Repo interface {
	Save(ctx context.Context, data *domain.RefreshToken) error
	Revoke(ctx context.Context, data *domain.RefreshToken) error
}

type ReadModel interface {
	GetPermissionById(ctx context.Context, id ulid.ULID) (PermissionList, error)
	Fetch(ctx context.Context) (RefreshTokenList, error)
	FindById(ctx context.Context, id ulid.ULID) (*domain.RefreshToken, error)
	FindByUserID(ctx context.Context, id ulid.ULID) (*domain.RefreshToken, error)
	FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
}
type PermissionList struct {
	Permissions []string `json:"data"`
	Count       int      `json:"count"`
}

var emptyPermissionList = PermissionList{
	Permissions: make([]string, 0),
	Count:       0,
}

// GetPermissionById implements ReadModel.
func (r *repo) GetPermissionById(ctx context.Context, id ulid.ULID) (PermissionList, error) {
	var itemCount int
	var urls []string
	rows, err := r.db.Query(
		ctx,
		`
			SELECT p.url, COUNT(p.url) AS url_count
			FROM account_roles ar
			JOIN role_permissions rp ON ar.role_id = rp.role_id
			JOIN permissions p ON rp.permission_id = p.id
			WHERE ar.account_id = $1
			GROUP BY p.url;
		`,
		id,
	)
	if err != nil {
		return emptyPermissionList, err
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		var count int
		if err := rows.Scan(&url, &count); err != nil {
			return emptyPermissionList, err
		}
		urls = append(urls, url)
		itemCount++
	}
	list := PermissionList{
		Permissions: urls,
		Count:       itemCount,
	}
	return list, nil
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{db: db}
}

func NewReadModel(db *pgxpool.Pool) ReadModel {
	return &repo{db: db}
}
