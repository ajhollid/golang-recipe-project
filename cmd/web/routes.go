package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/popnfresh234/recipe-app-golang/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)
	mux.Use(CorsMiddleware)
	mux.Get("/", handlers.Repo.Home)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)

	mux.Get("/user/signup", handlers.Repo.Signup)
	mux.Post("/user/signup", handlers.Repo.PostSignup)
	mux.Get("/user/logout", handlers.Repo.Logout)
	mux.Get("/recipe/details/{id}", handlers.Repo.RecipeDetails)

	mux.Route("/recipe", func(mux chi.Router) {
		//mux.Use(Auth)
		mux.Get("/new", handlers.Repo.NewRecipe)
		mux.Post("/new", handlers.Repo.PostNewRecipe)
		mux.Get("/edit/{id}", handlers.Repo.EditRecipe)
		mux.Post("/edit/{id}", handlers.Repo.PostEditRecipe)
		mux.Post("/delete/{id}", handlers.Repo.PostDeleteRecipe)
	})

	fileServer := http.FileServer(http.Dir("./web/static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
