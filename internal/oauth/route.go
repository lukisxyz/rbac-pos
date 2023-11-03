package oauth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"pos/domain"
	"pos/utils/httpresponse"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang-jwt/jwt/v5"
)

type accountRoute struct {
	svc          ServiceOAuth
	secret       string
	refreshToken uint
}

func NewRoute(
	svc ServiceOAuth,
	secret string,
	refreshToken uint,
) *accountRoute {
	return &accountRoute{
		svc:          svc,
		secret:       secret,
		refreshToken: refreshToken,
	}
}

func (p *accountRoute) Routes() *chi.Mux {
	r := chi.NewMux()
	r.Post("/login", p.login)
	r.Post("/logout", p.logout)
	r.Post("/request-token", p.requestaccesstoken)
	return r
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c loginRequest) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(8, 32)),
	)
}

func (p *accountRoute) logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	err := p.svc.Logout(ctx, body.Token)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	cookie := http.Cookie{
		Name:     "permissions",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	httpresponse.WriteMessage(w, http.StatusOK, "success logout")
}

func (p *accountRoute) login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	data, permissions, err := p.svc.Login(ctx, body.Email, body.Password)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	listPermission := strings.Join(permissions, ",")
	maxAge := p.refreshToken * 3600 * 24
	valuuEncrypted := base64.URLEncoding.EncodeToString([]byte(listPermission))

	cookie := http.Cookie{
		Name:     "permissions",
		Value:    valuuEncrypted,
		Path:     "/",
		MaxAge:   int(maxAge),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	httpresponse.WriteData(w, http.StatusOK, data, nil)
}

type tokenRequest struct {
	Token string `json:"token"`
}

func (h *accountRoute) requestaccesstoken(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	var body tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	reqToken := r.Header.Get("Authorization")
	splittedToken := strings.Split(reqToken, "Bearer ")
	if len(splittedToken) != 2 {
		httpresponse.WriteError(w, http.StatusUnauthorized, errors.New(http.StatusText((http.StatusUnauthorized))))
		ctx.Done()
		return
	}

	jwtToken := splittedToken[1]
	claims := &domain.Oauth{}
	token, err := jwt.ParseWithClaims(
		jwtToken,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(h.secret), nil
		},
	)

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		httpresponse.WriteError(w, http.StatusUnauthorized, err)
		ctx.Done()
		return
	}
	claims = token.Claims.(*domain.Oauth)

	accessToken, err := h.svc.RefreshToken(ctx, body.Token, claims.Id, claims.Email)
	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, err)
		ctx.Done()
		return
	}

	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, err)
		ctx.Done()
		return
	}

	httpresponse.WriteData(w, http.StatusOK, accessToken, nil)
}
