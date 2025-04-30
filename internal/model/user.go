package model

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
	Role     int    `json:"role" db:"role"`
}

// UserRole — тип роли пользователя (обычный пользователь или администратор)
// type UserRole string
type UserRole int

const (
	RoleUser  UserRole = 1
	RoleAdmin UserRole = 2
)

// РЕГИСТРАЦИЯ
type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Role     int    `json:"role"`
}

// ВХОД
type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// обновление данных (админка)
// type UpdateUserRequest struct {
// 	Username *string   `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
// 	Password *string   `json:"password,omitempty" validate:"omitempty,min=8"`
// 	Role     *UserRole `json:"role,omitempty" validate:"omitempty,oneof=user admin"`
// }

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// JWT-данные (хранятся в токене)
type TokenClaims struct {
	UserID int `json:"user_id"`
	Role   int `json:"role"`
	jwt.RegisteredClaims
}

// Права доступа (для middleware)
type Permission struct {
	Resource string   `json:"resource"` // Например, "movies", "actors"
	Actions  []string `json:"actions"`  // Например, ["create", "read"]
}
