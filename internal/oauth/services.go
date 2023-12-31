package oauth

import (
	"context"
	"errors"
	"pos/domain"
	"pos/internal/account"
	"pos/internal/permission"
	"pos/internal/role"
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

func (s *serviceOauth) RefreshToken(ctx context.Context, refreshToken string, uid ulid.ULID, email string) (accessToken string, err error) {
	_, err = s.readModel.FindByToken(ctx, refreshToken)
	if err != nil {
		return
	}
	jwtKey := []byte(s.secret)
	accessExpTime := time.Now().Add(time.Duration(s.accessExpTime) * time.Hour)
	claims := &domain.Oauth{
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
func (s *serviceOauth) Login(ctx context.Context, email, password string) (res *domain.LoginResponse, permissions []string, err error) {
	acc, err := s.accountReadModel.FindByEmail(ctx, email)
	if err != nil {
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(acc.Password),
		[]byte(password),
	); err != nil {
		return nil, permissions, ErrPasswordWrong
	}

	refreshExpTime := time.Now().Add(time.Duration(s.refreshExpTime) * 24 * time.Hour)

	jwtKey := []byte(s.secret)
	tokenRefreshString := utils.RandString(24)

	accessExpTime := time.Now().Add(time.Duration(s.accessExpTime) * time.Hour)
	claims := &domain.Oauth{
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

	oauthToken := domain.LoginResponse{
		AccessToken:  tokenAccessString,
		RefreshToken: tokenRefreshString,
		Type:         "Bearer",
		ExpiredAt:    refreshExpTime.Format(time.RFC3339),
		Scope:        "*",
	}

	refreshToken := domain.NewRefreshToken(
		acc.Id,
		tokenRefreshString,
		refreshExpTime,
	)

	data, err := s.readModel.GetPermissionById(ctx, acc.Id)
	if err != nil {
		return
	}
	permissions = data.Permissions

	_, err = s.readModel.FindByUserID(ctx, acc.Id)
	if errors.Is(err, ErrRefreshTokenNotFound) {
		err = s.repo.Save(ctx, &refreshToken)
		if err != nil {
			return
		}

		return &oauthToken, permissions, nil
	}
	err = ErrAlreadyLogin
	return
}

// Logout implements ServiceOAuth.
func (s *serviceOauth) Logout(ctx context.Context, token string) error {
	currentData, err := s.readModel.FindByToken(ctx, token)
	if err != nil {
		return err
	}
	return s.repo.Revoke(ctx, currentData)
}

type ServiceOAuth interface {
	Login(ctx context.Context, email, pwd string) (*domain.LoginResponse, []string, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string, uid ulid.ULID, email string) (accessToken string, err error)
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
