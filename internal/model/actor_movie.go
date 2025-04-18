package model

type ActorWithFilms struct {
	Actor Actor  `json:"actor"`
	Films []Film `json:"films"`
}

func ValidateGetActors() error {
	// TODO ... ДОБАВИТЬ ВАЛИДАЦИЮ, ЕСЛИ В ЭТОМ ЕСТЬ НЕОБХОДИМОСТЬ
	return nil
}
