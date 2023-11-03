package oauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type Oauth struct {
	Id    ulid.ULID `json:"ulid"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"Bearer"`
	ExpiredAt    string `json:"expired_at"`
	Scope        string `json:"scope"`
}

type RefreshToken struct {
	ID         ulid.ULID
	TokenValue string
	UserID     ulid.ULID
	CreatedAt  time.Time
	ExpiresAt  time.Time
	Revoked    bool
}

func NewRefreshToken(uid ulid.ULID, token string, expiredAt time.Time) RefreshToken {
	id := ulid.Make()

	res := RefreshToken{
		ID:         id,
		TokenValue: token,
		UserID:     uid,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiredAt,
		Revoked:    false,
	}

	return res
}
