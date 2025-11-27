package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Header'dan "Authorization" bilgisini al
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "Token yok!", 401)
				return
			}

			// "Bearer TOKEN" formatında mı diye bak
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Hatalı token formatı", 401)
				return
			}

			tokenStr := parts[1]

			// Token'ı çöz ve doğrula
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Geçersiz token", 401)
				return
			}

			// Her şey tamamsa içeri al
			next.ServeHTTP(w, r.WithContext(context.Background()))
		})
	}
}
