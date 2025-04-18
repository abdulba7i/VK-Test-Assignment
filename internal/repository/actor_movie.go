package repository

import (
	"context"
	"database/sql"
	"film-library/internal/model"
	"fmt"
)

type ActorMovieRepository interface {
	GetActorsWithFilms(ctx context.Context) (map[int]model.ActorWithFilms, error)
}

func NewActorMovieRepository(db *sql.DB) MovieRepository {
	return &Storage{
		db: db,
	}
}

func (s *Storage) GetActorsWithFilms(ctx context.Context) (map[int]model.ActorWithFilms, error) {
	const op = "storage.postgres.GetActorsWithFilms"

	query := `
        SELECT a.id, a.name, a.gender, a.date_of_birth, 
               f.id, f.name, f.description, f.release_date, f.rating
        FROM actors a
        JOIN actor_film af ON a.id = af.actor_id
        JOIN films f ON f.id = af.film_id`

	rows, err := s.db.QueryContext(ctx, query)
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
