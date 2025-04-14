package model

import (
	"errors"
	"time"
)

// Actor - основная сущность актёра
type Actor struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

// Validate - проверка корректности данных актёра
func (a *Actor) Validate() error {

	if a.Name == "" || len(a.Name) > 100 {
		return errors.New("имя актера не может быть пустым и не должно превышать 100 символов")
	}

	if a.Gender != "male" && a.Gender != "female" {
		return errors.New("несуществующий пол")
	}

	if a.DateOfBirth.After(time.Now()) {
		return errors.New("дата рождения не может быть в будущем")
	}

	if time.Now().Year()-a.DateOfBirth.Year() < 5 {
		return errors.New("актёр должен быть старше 5 лет")
	}

	return nil
}

// func NewActor(id int, name string, gender string, dateOfBirth time.Time) (*Actor, error) {
// 	actor := &Actor{
// 		Id:          id,
// 		Name:        name,
// 		Gender:      gender,
// 		DateOfBirth: dateOfBirth,
// 	}

// 	if err := actor.Validate(); err != nil {
// 		return nil, err
// 	}

// 	return actor, nil
// }
