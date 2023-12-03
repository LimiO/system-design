package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"onlinestore/pkg/web"
	"strings"
)

func (m *MiddlewareManager) Authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		parts := strings.Split(authorization, " ")
		if len(parts) == 2 {
			claims, err := m.tokenManager.GetClaims(parts[1])
			if err != nil {
				panic(fmt.Sprintf("failed to validate token: %v", err))
			}
			user, err := claims.GetSubject()
			if err != nil {
				panic(fmt.Sprintf("failed to get username from claims"))
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, web.UserIDKey{}, user)
			log.Println("authorized!")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`invalid token specified`))
	}
	return http.HandlerFunc(fn)
}
