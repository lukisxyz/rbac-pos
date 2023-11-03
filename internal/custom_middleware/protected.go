package custommiddleware

import (
	"encoding/base64"
	"errors"
	"net/http"
	"pos/utils/httpresponse"
	"strings"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func ProtectedMiddleware(grant string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				cookie, err := r.Cookie("permissions")
				if err != nil {
					httpresponse.WriteError(w, http.StatusUnauthorized, errors.New(http.StatusText((http.StatusUnauthorized))))
					ctx.Done()
					return
				}

				value, err := base64.URLEncoding.DecodeString(cookie.Value)
				if err != nil {
					httpresponse.WriteError(w, http.StatusUnauthorized, errors.New(http.StatusText((http.StatusUnauthorized))))
					ctx.Done()
					return
				}
				if err != nil {
					switch {
					case errors.Is(err, http.ErrNoCookie):
						httpresponse.WriteError(w, http.StatusBadRequest, errors.New(http.StatusText((http.StatusBadRequest))))
						ctx.Done()
						return
					case errors.Is(err, ErrInvalidValue):
						httpresponse.WriteError(w, http.StatusBadRequest, errors.New(http.StatusText((http.StatusBadRequest))))
						ctx.Done()
						return
					default:
						httpresponse.WriteError(w, http.StatusInternalServerError, errors.New(http.StatusText((http.StatusInternalServerError))))
						ctx.Done()
						return
					}
				}
				if strings.Contains(string(value), grant) {
					next.ServeHTTP(w, r)
				} else {
					httpresponse.WriteError(w, http.StatusUnauthorized, errors.New(http.StatusText((http.StatusUnauthorized))))
					ctx.Done()
					return
				}
			},
		)
	}
}
