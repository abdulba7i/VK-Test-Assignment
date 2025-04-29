package service

import (
	"context"
	"film-library/internal/model"
	"film-library/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	VerifyUser(ctx context.Context, username, password string) (string, *model.User, error)
	VerifyToken(tokenString string) (*model.TokenClaims, error)
}

type Actor interface {
	AddActor(ctx context.Context, actor model.Actor) error
	UpdateActor(ctx context.Context, actor model.Actor) error
	DeleteActor(ctx context.Context, id int) error
}

type Movie interface {
	AddMovie(ctx context.Context, film model.Film) error
	UpdateMovie(ctx context.Context, film model.Film) error
	DeleteMovie(ctx context.Context, id int) error
	GetFilms(ctx context.Context, sortBy string) ([]model.Film, error)
	SearchFilm(ctx context.Context, actor, film string) (model.Film, error)
}

type ActorMovie interface {
	GetAllActorWithFilms(ctx context.Context) (map[int]model.ActorWithFilms, error)
}

type Service struct {
	Authorization
	Actor
	Movie
	ActorMovie
}

func NewService(repos *repository.Repository, secret string) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, secret),
		Actor:         NewActorService(repos.Actor),
		Movie:         NewMovieService(repos.Movie),
		ActorMovie:    NewActorMovieService(repos.ActorMovie),
	}
}
