package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	// On panic, Recoverer sends back a 500 error - and brings things back up
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Post("/authenticate", app.authenticate)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)

	mux.Get("/movies", app.AllMovies)
	mux.Get("/movies/{id}", app.GetMovie)
	mux.Get("/genres", app.AllGenres)
	mux.Get("/movies/genres/{id}", app.AllMoviesByGenre)
	mux.Post("/graph", app.moviesGraphQL)

	// protected route
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authRequired)
		mux.Get("/movies", app.MovieCatalog)
		mux.Get("/movies/{id}", app.MovieForEdit)

		mux.Put("/movies/0", app.InsertMovie)
		mux.Patch("/movies/{id}", app.UpdateMovie)

		mux.Delete("/movies/{id}", app.DeleteMovie)
	})

	return mux
}
