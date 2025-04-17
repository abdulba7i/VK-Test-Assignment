package repository

import (
	"context"
	"database/sql"
	"film-library/internal/model"
	"fmt"
)

type MovieRepository interface {
	CreateFilm(ctx context.Context, film *model.Film) error
	UpdateFilm(ctx context.Context, film *model.Film) error
	DeleteFilm(ctx context.Context, id int) error
	GetAllFilms(ctx context.Context, sortBy string) ([]model.Film, error)
	SearchFilm(ctx context.Context, actor, film string) (model.Film, error)
	MovieExistsById(ctx context.Context, id int) (bool, error)
	MovieExistsByName(ctx context.Context, name string) (bool, error)
}

func NewMovieRepository(db *sql.DB) MovieRepository {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateFilm(ctx context.Context, film *model.Film) error {
	const op = "storage.postgres.AddedInfoFilm"

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	// Добавляем фильм
	var filmID int
	err = tx.QueryRowContext(ctx, `
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

func (s *Storage) UpdateFilm(ctx context.Context, film *model.Film) error {
	const op = "storage.postgres.ChangeInfoFilm"

	query := `UPDATE films SET name = $1, description = $2, release_date = $3, rating = $4 WHERE id = $5`

	_, err := s.db.ExecContext(ctx, query, film.Name, film.Description, film.Releasedate, film.Rating, film.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteFilm(ctx context.Context, id int) error {
	const op = "storage.postgres.DeleteInfoFilm"

	query := `DELETE FROM films WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAllFilms(ctx context.Context, sortBy string) ([]model.Film, error) {
	const op = "storage.postgres.GetAllFilms"

	orderClause := "ORDER BY rating DESC" // По умолчанию сортировка по рейтингу
	switch sortBy {
	case "name":
		orderClause = "ORDER BY name"
	case "release_date":
		orderClause = "ORDER BY release_date"
	}

	query := fmt.Sprintf(`
        SELECT f.id, f.name, f.description, f.release_date, f.rating, 
               a.id, a.name, a.gender, a.date_of_birth
        FROM films f
        JOIN actor_film af ON f.id = af.film_id
        JOIN actors a ON a.id = af.actor_id
        %s`, orderClause)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	filmMap := make(map[int]*model.Film)

	for rows.Next() {
		var film model.Film
		var actor model.Actor

		err = rows.Scan(
			&film.Id, &film.Name, &film.Description, &film.Releasedate, &film.Rating,
			&actor.Id, &actor.Name, &actor.Gender, &actor.DateOfBirth,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		// Если фильм ещё не добавлен в filmMap, добавляем его
		if _, exists := filmMap[film.Id]; !exists {
			film.ListActors = []model.Actor{} // Инициализируем пустой слайс актёров
			filmMap[film.Id] = &film
		}

		// Добавляем актёра к фильму
		filmMap[film.Id].ListActors = append(filmMap[film.Id].ListActors, actor)
	}

	// Преобразуем filmMap в слайс фильмов
	films := make([]model.Film, 0, len(filmMap))
	for _, film := range filmMap {
		films = append(films, *film)
	}

	return films, nil
}

func (s *Storage) SearchFilm(ctx context.Context, actor, film string) (model.Film, error) {
	const op = "storage.postgres.SearchFilm"

	// Ищем фильмы по фрагменту названия
	query := `
        SELECT f.id, f.name, f.description, f.release_date, f.rating, 
               a.id, a.name, a.gender, a.date_of_birth
        FROM films f
        JOIN actor_film af ON f.id = af.film_id
        JOIN actors a ON a.id = af.actor_id
        WHERE f.name ILIKE $1 AND a.name ILIKE $2`

	rows, err := s.db.QueryContext(ctx, query, "%"+film+"%", "%"+actor+"%")
	if err != nil {
		return model.Film{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var filmFound model.Film
	filmFound.ListActors = []model.Actor{} // Инициализация списка актёров

	for rows.Next() {
		var foundActor model.Actor

		// Если фильм ещё не найден, заполняем его данные
		if filmFound.Id == 0 {
			err = rows.Scan(
				&filmFound.Id, &filmFound.Name, &filmFound.Description,
				&filmFound.Releasedate, &filmFound.Rating,
				&foundActor.Id, &foundActor.Name, &foundActor.Gender, &foundActor.DateOfBirth,
			)
			if err != nil {
				return model.Film{}, fmt.Errorf("%s: %w", op, err)
			}
		} else {
			// Если фильм уже найден, сканируем только актёра
			err = rows.Scan(
				nil, nil, nil, nil, nil, // Пропускаем поля фильма
				&foundActor.Id, &foundActor.Name, &foundActor.Gender, &foundActor.DateOfBirth,
			)
			if err != nil {
				return model.Film{}, fmt.Errorf("%s: %w", op, err)
			}
		}

		// Добавляем актёра к фильму
		filmFound.ListActors = append(filmFound.ListActors, foundActor)
	}

	// Если фильм не найден, возвращаем ошибку
	if filmFound.Id == 0 {
		return model.Film{}, sql.ErrNoRows
	}

	return filmFound, nil
}

func (s *Storage) MovieExistsById(ctx context.Context, id int) (bool, error) {
	const op = "storage.postgres.ActorExistsById"

	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM films WHERE id = $1)`

	err := s.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (s *Storage) MovieExistsByName(ctx context.Context, name string) (bool, error) {
	const op = "storage.postgres.ActorExistsByName"

	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM films WHERE name = $1)`

	err := s.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}
