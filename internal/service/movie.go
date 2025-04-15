package service

import (
	"context"
	"errors"
	"film-library/internal/model"
	"film-library/internal/repository"
	"fmt"
)

type MovieService struct {
	repo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) AddMovie(ctx context.Context, film model.Film) error {
	exists, err := s.repo.MovieExistsByName(ctx, film.Name)
	if err != nil {
		return fmt.Errorf("Ошибка проверки фильма: %w", err)
	}
	if !exists {
		return errors.New("Фильм не найден")
	}
	return s.repo.CreateFilm(ctx, &film)
}

func (s *MovieService) UpdateMovie(ctx context.Context, film model.Film) error {
	exists, err := s.repo.MovieExistsById(ctx, film.Id)
	if err != nil {
		return fmt.Errorf("Ошибка проверки фильма: %w", err)
	}
	if !exists {
		return errors.New("Фильм не найден")
	}

	return s.repo.UpdateFilm(ctx, &film)
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int) error {
	exists, err := s.repo.MovieExistsById(ctx, id)
	if err != nil {
		return fmt.Errorf("Ошибка проверки фильма: %w", err)
	}
	if !exists {
		return errors.New("Фильм не найден")
	}

	return s.repo.DeleteFilm(ctx, id)
}
