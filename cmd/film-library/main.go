package main

import (
	"film-library/internal/config"
	"film-library/internal/handler"
	"film-library/internal/repository"
	"film-library/internal/service"
	slogpretty "film-library/internal/utils/handlers"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting film-library", slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	storage, err := repository.Connect(cfg.Database)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		fmt.Println("Storage connection error:", err)
		os.Exit(1)
	}

	actorService := service.NewActorService(storage)
	actorHandler := handler.NewActorHandler(*actorService)

	http.HandleFunc("/actors", actorHandler.HandleActors)
	http.HandleFunc("/actors/", actorHandler.HandleActor)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error("failed to start server", "error", err)
	}

	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog() // здесь преукрасили вывод логов для удобства
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("not supported env")
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

// actor1 := model.Actor{
// 	Name:        "Кама Пуля",
// 	Gender:      "female",
// 	DateOfBirth: "1111-11-11",
// }
// actor2 := model.Actor{
// 	Name:        "Рози Хантингтон-Уайтли",
// 	Gender:      "female",
// 	DateOfBirth: "1987-04-18",
// }
// actor3 := model.Actor{
// 	Name:        "Вин Дизель",
// 	Gender:      "male",
// 	DateOfBirth: "1967-07-18",
// }
// actor4 := model.Actor{
// 	Name:        "Мага Лезгин",
// 	Gender:      "male",
// 	DateOfBirth: "2222-12-22",
// }

// film := model.Film{
// 	Id:          1,
// 	Name:        "Пчеловод",
// 	Description: "Фантастическая история о компании неудачников, которые еще вчера жили в унылом мире и даже не подозревали, что однажды все перевернется с ног на голову",
// 	Releasedate: "2024-01-12",
// 	Rating:      8.3,
// 	ListActors: []model.Actor{
// 		actor1,
// 		actor4,
// 	},
// }
// film2 := model.Film{
// 	Name:        "Трансформеры 3: Тёмная сторона Луны",
// 	Description: "Сэм Уитвики должен снова спасти мир от десептиконов, которые на этот раз хотят использовать секретную технологию, спрятанную на Луне.",
// 	Releasedate: "2011-06-23",
// 	Rating:      9.2,
// 	ListActors: []model.Actor{
// 		actor2, // Рози Хантингтон-Уайтли
// 	},
// }
// film3 := model.Film{
// 	Id:          3,
// 	Name:        "Форсаж",
// 	Description: "Под прикрытием гонок Брайан О’Коннер должен внедриться в банду Доминика Торетто.",
// 	Releasedate: "2001-06-22",
// 	Rating:      8.5,
// 	ListActors: []model.Actor{
// 		actor3, // Вин Дизель
// 	},
// }
// film4 := model.Film{
// 	Id:          4,
// 	Name:        "Форсаж 5",
// 	Description: "Доминик Торрето и его команда, находясь в бегах, планируют последнее ограбление в Рио-де-Жанейро.",
// 	Releasedate: "2011-04-15",
// 	Rating:      9.0,
// 	ListActors: []model.Actor{
// 		actor3, // Вин Дизель
// 		actor1, // Джейсон Стэтхем
// 	},
// }

// err = storage.AddedInfoFilm(&film)
// if err != nil {
// 	log.Error("Ошибка при получении списка актёров", "error", err)
// 	os.Exit(1)
// }
// err = storage.AddedInfoFilm(&film2)
// if err != nil {
// 	log.Error("Ошибка при получении списка актёров", "error", err)
// 	os.Exit(1)
// }
// err = storage.AddedInfoFilm(&film3)
// if err != nil {
// 	log.Error("Ошибка при получении списка актёров", "error", err)
// 	os.Exit(1)
// }
// err = storage.AddedInfoFilm(&film4)
// if err != nil {
// 	log.Error("Ошибка при получении списка актёров", "error", err)
// 	os.Exit(1)
// }

// actorsWithFilms, err := storage.GetActorsWithFilms()
// if err != nil {
// 	log.Error("Ошибка при получении списка актёров", "error", err)
// 	os.Exit(1)
// }
// for _, actorWithFilms := range actorsWithFilms {
// 	log.Info("Актёр", "actor", actorWithFilms.Actor)
// 	for _, film := range actorWithFilms.Films {
// 		log.Info("  Фильм", "film", film)
// 	}
// }
