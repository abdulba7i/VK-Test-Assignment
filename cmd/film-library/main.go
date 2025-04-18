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

	// Actor

	actorService := service.NewActorService(storage)
	actorHandler := handler.NewActorHandler(*actorService)

	// Movie

	movieService := service.NewMovieService(storage)
	movieHandler := handler.NewMovieHandler(*movieService)

	// ActorMovie

	actormovieService := service.NewActorMovieService(storage)
	actormovieHandler := handler.NewActorMovieHandler(*actormovieService)

	//

	http.HandleFunc("/actors", actorHandler.HandleActorPost)
	http.HandleFunc("/actor/", actorHandler.HandleActorPut)

	//

	http.HandleFunc("/films", movieHandler.HandleMoviePost)
	http.HandleFunc("/film/", movieHandler.HandleMoviePut)
	http.HandleFunc("/films_get_list", movieHandler.GetAllFilms) // GET /films
	http.HandleFunc("/films/search", movieHandler.SearchFilm)    // GET /films/search

	//

	http.HandleFunc("/actors_films", actormovieHandler.HandleActorMovieGet)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error("failed to start server", "error", err)
	}

	_ = storage
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
