package service

import (
	"context"
	"film-library/internal/model"
	"film-library/internal/repository"
	"fmt"
)

type ActorMovieService struct {
	repo repository.ActorMovieRepository
}

func NewActorMovieService(repo repository.ActorMovieRepository) *ActorMovieService {
	return &ActorMovieService{repo: repo}
}

func (s *ActorMovieService) GelAllActorWithFilms(ctx context.Context) (map[int]model.ActorWithFilms, error) {
	ListActors, err := s.repo.GetActorsWithFilms(ctx)
	if err != nil {
		return map[int]model.ActorWithFilms{}, fmt.Errorf("Ошибка получения списка всех актеров: %w", err)
	}

	if len(ListActors) == 0 {
		return map[int]model.ActorWithFilms{}, nil
	}

	return ListActors, nil
}
