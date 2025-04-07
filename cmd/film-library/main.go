package main

import (
	"film-library/internal/config"
	"film-library/internal/repository"
	slogpretty "film-library/internal/utils/handlers"
	"fmt"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
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

	actorsWithFilms, err := storage.GetActorsWithFilms()
	if err != nil {
		log.Error("Ошибка при получении списка актёров", "error", err)
		os.Exit(1)
	}

	for _, actorWithFilms := range actorsWithFilms {
		log.Info("Актёр", "actor", actorWithFilms.Actor)
		for _, film := range actorWithFilms.Films {
			log.Info("  Фильм", "film", film)
		}
	}

	fmt.Println(actorsWithFilms)

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
