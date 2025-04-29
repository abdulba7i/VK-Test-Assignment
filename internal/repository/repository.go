package repository

import (
	"context"
	"database/sql"
	"film-library/internal/model"
)

// AuthRepository
type Authorization interface {
	CreateUser(ctx context.Context, user *model.User) error
	VerifyUser(ctx context.Context, username string) (*model.User, error)
}

// ActorRepository
type Actor interface {
	CreateActor(ctx context.Context, actor *model.Actor) error
	UpdateActor(ctx context.Context, actor *model.Actor) error
	DeleteActor(ctx context.Context, id int) error
	ActorExistsById(ctx context.Context, id int) (bool, error)
	ActorExistsByName(ctx context.Context, name string) (bool, error)
}

// MovieRepository
type Movie interface {
	CreateFilm(ctx context.Context, film *model.Film) error
	UpdateFilm(ctx context.Context, film *model.Film) error
	DeleteFilm(ctx context.Context, id int) error
	GetAllFilms(ctx context.Context, sortBy string) ([]model.Film, error)   // получение списка фильмов
	SearchFilm(ctx context.Context, actor, film string) (model.Film, error) // поиск фильмов
	MovieExistsById(ctx context.Context, id int) (bool, error)
	MovieExistsByName(ctx context.Context, name string) (bool, error)
}

// ActorMovieRepository
type ActorMovie interface {
	GetActorsWithFilms(ctx context.Context) (map[int]model.ActorWithFilms, error)
}

type Repository struct {
	Authorization
	Actor
	Movie
	ActorMovie
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		Actor:         NewActorRepository(db),
		Movie:         NewMovieRepository(db),
		ActorMovie:    NewActorMovieRepository(db),
	}
}
