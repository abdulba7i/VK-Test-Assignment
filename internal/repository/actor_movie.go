package repository

import (
	"database/sql"
	"film-library/internal/model"
	"fmt"
)

func (s *Storage) GetAllFilms(sortBy string) ([]model.Film, error) {
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

	rows, err := s.db.Query(query)
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

func (s *Storage) SearchFilm(actor, film string) (model.Film, error) {
	const op = "storage.postgres.SearchFilm"

	// Ищем фильмы по фрагменту названия
	query := `
        SELECT f.id, f.name, f.description, f.release_date, f.rating, 
               a.id, a.name, a.gender, a.date_of_birth
        FROM films f
        JOIN actor_film af ON f.id = af.film_id
        JOIN actors a ON a.id = af.actor_id
        WHERE f.name ILIKE $1 AND a.name ILIKE $2`

	rows, err := s.db.Query(query, "%"+film+"%", "%"+actor+"%")
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

func (s *Storage) GetActorsWithFilms() (map[int]model.ActorWithFilms, error) {
	const op = "storage.postgres.GetActorsWithFilms"

	query := `
        SELECT a.id, a.name, a.gender, a.date_of_birth, 
               f.id, f.name, f.description, f.release_date, f.rating
        FROM actors a
        JOIN actor_film af ON a.id = af.actor_id
        JOIN films f ON f.id = af.film_id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	actorsWithFilms := make(map[int]model.ActorWithFilms)

	for rows.Next() {
		var actor model.Actor
		var film model.Film

		err = rows.Scan(
			&actor.Id, &actor.Name, &actor.Gender, &actor.DateOfBirth,
			&film.Id, &film.Name, &film.Description, &film.Releasedate, &film.Rating,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		// Если актёр ещё не добавлен в мапу, добавляем его
		if _, exists := actorsWithFilms[actor.Id]; !exists {
			actorsWithFilms[actor.Id] = model.ActorWithFilms{
				Actor: actor,
				Films: []model.Film{},
			}
		}

		// Добавляем фильм к актёру
		actorWithFilms := actorsWithFilms[actor.Id]
		actorWithFilms.Films = append(actorWithFilms.Films, film)
		actorsWithFilms[actor.Id] = actorWithFilms
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return actorsWithFilms, nil
}
