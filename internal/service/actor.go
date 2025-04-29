package service

import (
	"context"
	"errors"
	"film-library/internal/model"
	"film-library/internal/repository"
	"fmt"
)

type ActorService struct {
	// repdo repository.ActorRepository
	repo repository.Actor
}

func NewActorService(repo repository.Actor) *ActorService {
	return &ActorService{repo: repo}
}

func (s *ActorService) AddActor(ctx context.Context, actor model.Actor) error {
	exists, err := s.repo.ActorExistsByName(ctx, actor.Name)
	if err != nil {
		return fmt.Errorf("ошибка проверки актёра: %w", err)
	}
	if !exists {
		return errors.New("актёр не найден")
	}

	return s.repo.CreateActor(ctx, &actor)
}

func (s *ActorService) UpdateActor(ctx context.Context, actor model.Actor) error {
	exists, err := s.repo.ActorExistsById(ctx, actor.Id)
	if err != nil {
		return fmt.Errorf("ошибка проверки актёра: %w", err)
	}
	if !exists {
		return errors.New("актёр не найден")
	}

	if err := actor.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации актёра: %w", err)
	}

	return s.repo.UpdateActor(ctx, &actor)
}

func (s *ActorService) DeleteActor(ctx context.Context, id int) error {
	// TODO: ... могу ли я удалять актера, есть он привязан к какому-либо фильму??
	return s.repo.DeleteActor(ctx, id)
}
