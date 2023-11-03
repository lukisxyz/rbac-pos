package custommiddleware

import (
	"context"
	"errors"
	"net/http"
	"pos/domain"
	"pos/utils/httpresponse"
	"pos/utils/key"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = ""

func SetJwtSecret(j string) {
	jwtSecret = j
}

func AuthJwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqToken := r.Header.Get("Authorization")
			splittedToken := strings.Split(reqToken, "Bearer ")
			ctx := r.Context()
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
					return []byte(jwtSecret), nil
				},
			)
			if err != nil {
				httpresponse.WriteError(w, http.StatusUnauthorized, err)
				ctx.Done()
				return
			}
			if !token.Valid {
				httpresponse.WriteError(w, http.StatusUnauthorized, errors.New(http.StatusText((http.StatusUnauthorized))))
				ctx.Done()
				return
			}

			c := context.WithValue(
				ctx,
				key.UserValueKey,
				claims,
			)
			next.ServeHTTP(w, r.WithContext(c))
		},
	)
}
