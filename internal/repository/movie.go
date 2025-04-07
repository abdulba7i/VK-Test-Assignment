package repository

import (
	"database/sql"
	"film-library/internal/model"
	"fmt"
)

func (s *Storage) AddedInfoFilm(film *model.Film) error {
	const op = "storage.postgres.AddedInfoFilm"

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	// Добавляем фильм
	var filmID int
	err = tx.QueryRow(`
        INSERT INTO films (name, description, release_date, rating)
        VALUES ($1, $2, $3, $4)
        RETURNING id`,
		film.Name, film.Description, film.Releasedate, film.Rating,
	).Scan(&filmID)
	if err != nil {
		return fmt.Errorf("%s: failed to insert film: %w", op, err)
	}

	// Добавляем актёров
	for _, actor := range film.ListActors {
		var actorID int

		// Проверяем, существует ли актёр
		err := tx.QueryRow(`
            SELECT id FROM actors WHERE name = $1`,
			actor.Name,
		).Scan(&actorID)

		if err != nil {
			if err == sql.ErrNoRows {
				// Актёр не существует, создаём нового
				err = tx.QueryRow(`
                    INSERT INTO actors (name, gender, date_of_birth)
                    VALUES ($1, $2, $3)
                    RETURNING id`,
					actor.Name, actor.Gender, actor.DateOfBirth,
				).Scan(&actorID)
				if err != nil {
					return fmt.Errorf("%s: failed to insert actor: %w", op, err)
				}
			} else {
				return fmt.Errorf("%s: failed to check actor existence: %w", op, err)
			}
		}

		// Связываем фильм и актёра
		_, err = tx.Exec(`
            INSERT INTO actor_film (film_id, actor_id)
            VALUES ($1, $2)`,
			filmID, actorID,
		)
		if err != nil {
			return fmt.Errorf("%s: failed to insert film-actor link: %w", op, err)
		}
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateFilm(film *model.Film) error {
	const op = "storage.postgres.ChangeInfoFilm"

	query := `UPDATE films SET name = $1, description = $2, release_date = $3, rating = $4 WHERE id = $5`

	_, err := s.db.Exec(query, film.Name, film.Description, film.Releasedate, film.Rating, film.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteInfoFilm(id int) error {
	const op = "storage.postgres.DeleteInfoFilm"

	query := `DELETE FROM films WHERE id = $1`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
