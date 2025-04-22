package repository

import (
	"context"
	"database/sql"
	"errors"
	"film-library/internal/model"
	"fmt"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	VerifyUser(ctx context.Context, username string) (*model.User, error)
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(ctx context.Context, user *model.User) error {
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO users (name, password, role_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, user.Username, user.Password, user.Role,
	).Scan(&user.ID)

	return err
}

func (s *Storage) VerifyUser(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, password, role_id
		FROM users
		WHERE name = $1
	`, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}
