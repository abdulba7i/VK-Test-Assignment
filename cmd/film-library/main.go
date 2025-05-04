package main

import (
	"film-library/internal/config"
	"film-library/internal/handler"
	"film-library/internal/middleware"
	"film-library/internal/repository"
	"film-library/internal/service"
	slogpretty "film-library/internal/utils/handlers"
	"fmt"
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
	actorHandler := handler.NewActorHandler(actorService)

	// Movie

	movieService := service.NewMovieService(storage)
	movieHandler := handler.NewMovieHandler(movieService)

	// ActorMovie

	actormovieService := service.NewActorMovieService(storage)
	actormovieHandler := handler.NewActorMovieHandler(actormovieService)

	// Auth
	secret := os.Getenv("SECRET_KEY")
	authService := service.NewAuthService(storage, secret)
	authHandler := handler.NewAuthHandler(authService)

	//

	http.HandleFunc("/actors", middleware.RequireAuth([]byte(secret))(actorHandler.HandleActorPost))       // create actor
	http.HandleFunc("/actor/", middleware.RequireAuth([]byte(secret))(actorHandler.HandleActorPut))        // update actor
	http.HandleFunc("/actor_del/", middleware.RequireAuth([]byte(secret))(actorHandler.HandleActorDelete)) // update actor

	//

	http.HandleFunc("/films", middleware.RequireAuth([]byte(secret))(movieHandler.HandleMoviePost))
	http.HandleFunc("/film/", middleware.RequireAuth([]byte(secret))(movieHandler.HandleMoviePut))
	http.HandleFunc("/film_del/", middleware.RequireAuth([]byte(secret))(movieHandler.HandleMovieDelete))

	http.HandleFunc("/films_get_list/", middleware.RequireAuth([]byte(secret))(movieHandler.GetAllFilms)) // GET /films

	http.HandleFunc("/films/search", middleware.RequireAuth([]byte(secret))(movieHandler.SearchFilm)) // GET /films/search

	//

	http.HandleFunc("/actors_films", middleware.RequireAuth([]byte(secret))(actormovieHandler.HandleActorMovieGet))

	//

	http.HandleFunc("/sign_up", authHandler.HandleAuthPost)
	http.HandleFunc("/sign_in", authHandler.HandleAuthPost)

	//

	// http.HandleFunc("/user-data", middleware.RequireRole([]byte(secret), int(model.RoleAdmin))(UserDataHandler))
	// http.HandleFunc("/admin-only", middleware.RequireRole([]byte(secret), int(model.RoleUser))(AdminHandler))

	//

	go func() {
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Error("failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// err = http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Error("failed to start server", "error", err)
	// }

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
