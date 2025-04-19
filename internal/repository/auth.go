package repository

import (
	"database/sql"
	"errors"
	"film-library/internal/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) CreateUser(user *model.User) error {
	err := r.db.QueryRow(`
	INSERT INTO users (name, password, role_id)
	VALUES ($1, $2, $3)
	RETURNING id`,
		user.Username, user.Password, user.Role,
	).Scan(&user.ID)

	return err
}

func (r *AuthRepository) VerifyUser(username, password string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(`
	SELECT id, name, password, role_id
	FROM users
	WHERE name = $1 AND password = $2
	`, username, password).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
