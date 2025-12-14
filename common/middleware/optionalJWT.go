package middleware

import (
	"context"
	"fungo/common/jwts"
	"net/http"
	"strings"
)

func OptionalJWT(secret string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if strings.HasPrefix(strings.TrimSpace(auth), "Bearer ") {
				tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
				claims, err := jwts.ParseToken(tokenStr, secret)
				if err == nil {
					ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
