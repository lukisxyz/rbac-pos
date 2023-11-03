package custommiddleware

import (
	"fmt"
	"net/http"
)

func ProtectedMiddleware(grant string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(grant)
				next.ServeHTTP(w, r)
			},
		)
	}
}
