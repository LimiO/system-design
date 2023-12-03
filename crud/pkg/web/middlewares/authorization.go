package middlewares

import (
	"net/http"
	"strings"
)

func Authorize(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		parts := strings.Split(authorization, " ")
		if len(parts) == 2 {
			// TODO(albert-si) check valid token
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`invalid token specified`))
	}
	return http.HandlerFunc(fn)
}
