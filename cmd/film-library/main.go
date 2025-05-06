package main

import (
	"film-library/internal/config"
	"film-library/internal/handler"
	"film-library/internal/repository"
	"film-library/internal/service"
	slogpretty "film-library/internal/utils/handlers"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
	version  = "1.0.0"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	secret := os.Getenv("SECRET_KEY")
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting film-library",
		slog.String("env", cfg.Env),
		slog.String("version", version),
	)

	storage, err := repository.Connect(cfg.Database)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	repositories := repository.NewRepository(storage.DB())
	services := service.NewService(repositories, secret)
	router := handler.InitRoute(services)

	log.Info("Server is running on port :8080")
	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Error("server exited with error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gracefully...")
}

// http://localhost:8080/films/search?actor=Рози Хантингтон-Уайтли&movie=Трансформеры 3: Тёмная сторона Луны
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
