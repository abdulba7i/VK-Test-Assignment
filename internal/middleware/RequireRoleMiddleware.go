package middleware

import (
	"film-library/internal/utils/response"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func RequireRole(secretKey []byte, requiredRole int) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteJSONError(w, "Missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})
			if err != nil || !token.Valid {
				response.WriteJSONError(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				response.WriteJSONError(w, "Invalid claims", http.StatusUnauthorized)
				return
			}

			roleFloat, ok := claims["role"].(float64)
			if !ok {
				response.WriteJSONError(w, "Role not found", http.StatusForbidden)
				return
			}
			role := int(roleFloat)

			if role != requiredRole {
				response.WriteJSONError(w, "Forbidden", http.StatusForbidden)
				return
			}

			next(w, r)
		}
	}
}
