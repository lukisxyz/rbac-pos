package oauth

import (
	"context"
	"errors"
	"pos/account"
	"pos/permission"
	"pos/role"
	"pos/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordWrong = errors.New("login: wrong password")
	ErrAlreadyLogin  = errors.New("login: already login")
)

type serviceOauth struct {
	roleReadModel       role.ReadModel
	permissionReadModel permission.ReadModel
	accountReadModel    account.ReadModel
	repo                Repo
	readModel           ReadModel
	secret              string
	refreshExpTime      uint
	accessExpTime       uint
}

func (s *serviceOauth) RefreshToken(ctx context.Context, uid ulid.ULID, name, email string) (accessToken string, err error) {
	_, err = s.readModel.FindByUserID(ctx, uid)
	if err != nil {
		return
	}
	jwtKey := []byte(s.secret)
	accessExpTime := time.Now().Add(time.Duration(s.accessExpTime) * time.Second)
	claims := &Oauth{
		Id:    uid,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpTime),
		},
	}
	tokenAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenAccess.SignedString(jwtKey)
}

// Login implements ServiceOAuth.
func (s *serviceOauth) Login(ctx context.Context, email, password string) (res *LoginResponse, err error) {
	acc, err := s.accountReadModel.FindByEmail(ctx, email)
	if err != nil {
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(acc.Password),
		[]byte(password),
	); err != nil {
		return nil, ErrPasswordWrong
	}

	refreshExpTime := time.Now().Add(time.Duration(s.refreshExpTime) * 24 * time.Second)

	jwtKey := []byte(s.secret)
	tokenRefreshString := utils.RandString(24)

	accessExpTime := time.Now().Add(time.Duration(s.accessExpTime) * time.Second)
	claims := &Oauth{
		Id:    acc.Id,
		Email: acc.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpTime),
		},
	}
	tokenAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenAccessString, errSign := tokenAccess.SignedString(jwtKey)
	if errSign != nil {
		err = errSign
		return
	}

	oauthToken := LoginResponse{
		AccessToken:  tokenAccessString,
		RefreshToken: tokenRefreshString,
		Type:         "Bearer",
		ExpiredAt:    refreshExpTime.Format(time.RFC3339),
		Scope:        "*",
	}

	refreshToken := NewRefreshToken(
		acc.Id,
		tokenRefreshString,
		refreshExpTime,
	)

	_, err = s.readModel.FindByUserID(ctx, acc.Id)
	if errors.Is(err, ErrRefreshTokenNotFound) {
		err = s.repo.Save(ctx, &refreshToken)
		if err != nil {
			return
		}

		return &oauthToken, nil
	}
	err = ErrAlreadyLogin
	return
}

// Logout implements ServiceOAuth.
func (s *serviceOauth) Logout(ctx context.Context, id ulid.ULID) error {
	currentData, err := s.readModel.FindByUserID(ctx, id)
	if err != nil {
		return err
	}
	currentData.Revoked = true
	return s.repo.Save(ctx, currentData)
}

type ServiceOAuth interface {
	Login(ctx context.Context, email, pwd string) (*LoginResponse, error)
	Logout(ctx context.Context, id ulid.ULID) error
	RefreshToken(ctx context.Context, uid ulid.ULID, name, email string) (accessToken string, err error)
}

func NewServiceOAuth(
	accountReadModel account.ReadModel,
	repo Repo,
	readModel ReadModel,
	secret string,
	refreshExpTime uint,
	accessExpTime uint,
	roleReadModel role.ReadModel,
	permissionReadModel permission.ReadModel,
) ServiceOAuth {
	return &serviceOauth{
		accountReadModel:    accountReadModel,
		repo:                repo,
		readModel:           readModel,
		secret:              secret,
		refreshExpTime:      refreshExpTime,
		accessExpTime:       accessExpTime,
		roleReadModel:       roleReadModel,
		permissionReadModel: permissionReadModel,
	}
}
