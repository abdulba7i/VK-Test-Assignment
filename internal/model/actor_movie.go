package model

type ActorWithFilms struct {
	Actor Actor  `json:"actor"`
	Films []Film `json:"films"`
}

func (a *ActorWithFilms) ValidateGetActors() error {
	// TODO ... ДОБАВИТЬ ВАЛИДАЦИЮ, ЕСЛИ В ЭТОМ ЕСТЬ НЕОБХОДИМОСТЬ
	return nil
}
