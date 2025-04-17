package model

import (
	"fmt"
	"time"
)

type Film struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Releasedate time.Time `json:"release_date"`
	Rating      float32   `json:"rating"`
	ListActors  []Actor   `json:"list_actors"`
}

// Validate - проверка данных фильма
func (f *Film) Validate() error {
	if f.Description == "" || len(f.Description) > 150 {
		return fmt.Errorf("название фильма должно быть от 1 до 150 символов")
	}
	if len(f.Description) > 1000 {
		return fmt.Errorf("описание не должно превышать 1000 символов")
	}
	if f.Rating < 0 || f.Rating > 10 {
		return fmt.Errorf("рейтинг должен быть от 0 до 10")
	}
	return nil
}

func ValidateSortFilm(sortBy string) error {
	if sortBy != "name" && sortBy != "release_date" && sortBy != "" {
		return fmt.Errorf("Некорректная сортировка")
	}
	return nil
}

func (a *Film) ValidateFilmSearchParams(filmName, actorName string) error {
	if filmName == "" && actorName == "" {
		return fmt.Errorf("необходимо указать либо название фильма, либо имя актёра")
	}

	if filmName != "" {
		if len(filmName) > 150 {
			return fmt.Errorf("название фильма слишком длинное (макс. 150 символов)")
		}
		if len(filmName) < 2 {
			return fmt.Errorf("название фильма слишком короткое (мин. 2 символа)")
		}
	}

	if actorName != "" {
		if len(actorName) > 100 {
			return fmt.Errorf("имя актёра слишком длинное (макс. 100 символов)")
		}
		if len(actorName) < 2 {
			return fmt.Errorf("имя актёра слишком короткое (мин. 2 символа)")
		}
	}

	return nil
}
