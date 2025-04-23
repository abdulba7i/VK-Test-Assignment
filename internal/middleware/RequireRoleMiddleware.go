package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func RequireRole(secretKey []byte, requiredRole int) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid claims", http.StatusUnauthorized)
				return
			}

			roleFloat, ok := claims["role"].(float64)
			if !ok {
				http.Error(w, "Role not found", http.StatusForbidden)
				return
			}
			role := int(roleFloat)

			if role != requiredRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next(w, r)
		}
	}
}
