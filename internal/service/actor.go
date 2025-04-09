package service

import (
	"context"
	"errors"
	"film-library/internal/model"
	"film-library/internal/repository"
	"fmt"
	"time"
)

type ActorService struct {
	repo repository.ActorRepository
}

func (s *ActorService) AddActor(ctx context.Context, actor model.Actor) error {
	if actor.Name == "" || len(actor.Name) > 100 {
		return errors.New("Имя актера не может быть пустым и не должно превышать 100 символов")
	}

	if actor.Gender != "male" && actor.Gender != "female" {
		return errors.New("Не существующий пол")
	}

	if actor.DateOfBirth.After(time.Now()) {
		return errors.New("Дата рождения не может быть в будущем")
	}

	if time.Now().Year()-actor.DateOfBirth.Year() < 5 {
		return errors.New("актёр должен быть старше 5 лет")
	}

	return s.repo.AddedInfoActor(ctx, &actor)
}

func (s *ActorService) UpdateActor(ctx context.Context, actor model.Actor) error {
	exists, err := s.repo.ActorExistsById(ctx, actor.Id)
	if err != nil {
		return fmt.Errorf("ошибка проверки актёра: %w", err)
	}
	if !exists {
		return errors.New("актёр не найден")
	}

	if actor.Name == "" || len(actor.Name) > 100 {
		return errors.New("Имя актера не может быть пустым и не должно превышать 100 символов")
	}

	if !actor.DateOfBirth.IsZero() && actor.DateOfBirth.After(time.Now()) {
		return errors.New("дата рождения не может быть в будущем")
	}

	return s.repo.UpdateActor(ctx, &actor)
}

func (s *ActorService) DeleteInfoActor(ctx context.Context, id int) error {
	exists, err := s.repo.ActorExistsById(ctx, id)
	if err != nil {
		return fmt.Errorf("ошибка проверки актёра: %w", err)
	}
	if !exists {
		return errors.New("актёр не найден")
	}

	// TODO: ... могу ли я удалять актера, есть он привязан к какому-либо фильму??
	return s.repo.DeleteInfoActor(ctx, id)
}
