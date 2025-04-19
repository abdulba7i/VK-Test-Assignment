package model

import "time"

type User struct {
	ID       int      `json:"id" db:"id"`
	Username string   `json:"username" db:"username"`
	Password string   `json:"-" db:"password"` // Пароль исключён из JSON (не возвращаем клиенту)
	Role     UserRole `json:"role" db:"role"`
	// CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserRole — тип роли пользователя (обычный пользователь или администратор)
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// регистрация
type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

// Заход
type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// обновление данных (админка)
type UpdateUserRequest struct {
	Username *string   `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Password *string   `json:"password,omitempty" validate:"omitempty,min=8"`
	Role     *UserRole `json:"role,omitempty" validate:"omitempty,oneof=user admin"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// JWT-данные (хранятся в токене)
type TokenClaims struct {
	UserID int      `json:"user_id"`
	Role   UserRole `json:"role"`
}

// Права доступа (для middleware)
type Permission struct {
	Resource string   `json:"resource"` // Например, "movies", "actors"
	Actions  []string `json:"actions"`  // Например, ["create", "read"]
}
