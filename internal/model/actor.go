package model

import (
	"fmt"
)

// Actor - основная сущность актёра
type Actor struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

// Validate - проверка корректности данных актёра
func (a *Actor) Validate() error {
	if a.Name == "" || len(a.Name) > 100 {
		return fmt.Errorf("имя актёра должно быть от 1 до 100 символов")
	}
	if a.Gender != "male" && a.Gender != "female" {
		return fmt.Errorf("пол должен быть 'male' или 'female'")
	}
	return nil
}
