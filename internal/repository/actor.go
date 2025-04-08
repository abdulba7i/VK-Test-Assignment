package repository

import (
	"context"
	"film-library/internal/model"
	"fmt"
)

func (s *Storage) AddedInfoActor(ctx context.Context, actor *model.Actor) error {
	const op = "storage.postgres.AddedInfoActor"

	query := `INSERT INTO actors (name, gender, date_of_birth) VALUES ($1, $2, $3) RETURNING id`
	err := s.db.QueryRowContext(ctx, query, actor.Name, actor.Gender, actor.DateOfBirth).Scan(&actor.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) UpdateActor(ctx context.Context, actor *model.Actor) error {
	const op = "storage.postgres.ChangeInfoActor"

	query := `UPDATE actors SET name = $1, gender = $2, date_of_birth = $3 WHERE id = $4`
	// Используем ExecContext вместо Exec
	_, err := s.db.ExecContext(ctx, query, actor.Name, actor.Gender, actor.DateOfBirth, actor.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteInfoActor(ctx context.Context, id int) error {
	const op = "storage.postgres.DeleteInfoActor"

	query := `DELETE FROM actors WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
