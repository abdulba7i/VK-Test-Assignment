package utils

import "net/http"

func IsAdmin(r *http.Request) bool {
	roleVal := r.Context().Value("role")
	if roleVal == nil {
		return false
	}
	role, ok := roleVal.(int)
	return ok && role == 1
}
