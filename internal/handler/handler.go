package handler

import (
	"film-library/internal/middleware"
	"film-library/internal/service"
	"net/http"
	"os"
)

func InitRoute(services *service.Service) *http.ServeMux {
	mux := http.NewServeMux()

	secret := os.Getenv("SECRET_KEY")

	actorHandler := NewActorHandler(services.Actor)
	movieHandler := NewMovieHandler(services.Movie)
	actormovieHandler := NewActorMovieHandler(services.ActorMovie)
	authHandler := NewAuthHandler(services.Authorization)

	mux.HandleFunc("/actors", middleware.RequireAuth([]byte(secret))(actorHandler.HandleActorPost))
	mux.HandleFunc("/actor/", middleware.RequireAuth([]byte(secret))(actorHandler.HandleActorPut))
	mux.HandleFunc("/actor_del/", middleware.RequireAuth([]byte(secret))(actorHandler.HandleActorDelete))

	mux.HandleFunc("/films", middleware.RequireAuth([]byte(secret))(movieHandler.HandleMoviePost))
	mux.HandleFunc("/film/", middleware.RequireAuth([]byte(secret))(movieHandler.HandleMoviePut))
	mux.HandleFunc("/film_del/", middleware.RequireAuth([]byte(secret))(movieHandler.HandleMovieDelete))
	mux.HandleFunc("/films_get_list/", middleware.RequireAuth([]byte(secret))(movieHandler.GetAllFilms))
	mux.HandleFunc("/films/search", middleware.RequireAuth([]byte(secret))(movieHandler.SearchFilm))

	// Актёры + фильмы
	mux.HandleFunc("/actors_films", middleware.RequireAuth([]byte(secret))(actormovieHandler.HandleActorMovieGet))

	// Аутентификация
	mux.HandleFunc("/sign_up", authHandler.HandleAuthPost)
	mux.HandleFunc("/sign_in", authHandler.HandleAuthPost)

	return mux
}
