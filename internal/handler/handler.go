package handler

import (
	"film-library/internal/middleware"
	"film-library/internal/service"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "film-library/docs"
)

func InitRoute(services *service.Service) *http.ServeMux {
	mux := http.NewServeMux()
	secret := os.Getenv("SECRET_KEY")

	actorHandler := NewActorHandler(services.Actor)
	movieHandler := NewMovieHandler(services.Movie)
	actormovieHandler := NewActorMovieHandler(services.ActorMovie)
	authHandler := NewAuthHandler(services.Authorization)

	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Актеры
	mux.HandleFunc("/actor_create", middleware.RequireAuth([]byte(secret))(actorHandler.HandleActorPost))
	mux.HandleFunc("/actor_update", middleware.RequireAuth([]byte(secret))(actorHandler.HandleActorPut))
	mux.HandleFunc("/actor_delete/", middleware.RequireAuth([]byte(secret))(actorHandler.HandleActorDelete))

	// Фильмы
	mux.HandleFunc("/film_create", middleware.RequireAuth([]byte(secret))(movieHandler.HandleMoviePost))
	mux.HandleFunc("/film_update", middleware.RequireAuth([]byte(secret))(movieHandler.HandleMoviePut))
	mux.HandleFunc("/film_delete/", middleware.RequireAuth([]byte(secret))(movieHandler.HandleMovieDelete))
	mux.HandleFunc("/films_get_list", middleware.RequireAuth([]byte(secret))(movieHandler.GetAllFilms))
	mux.HandleFunc("/films/search", middleware.RequireAuth([]byte(secret))(movieHandler.SearchFilm))

	// Актёры + фильмы
	mux.HandleFunc("/get_list_actors_films", middleware.RequireAuth([]byte(secret))(actormovieHandler.HandleActorMovieGet))

	// Аутентификация
	mux.HandleFunc("/auth/sign_up", authHandler.HandleAuthPost)
	mux.HandleFunc("/auth/sign_in", authHandler.HandleAuthPost)

	return mux
}
