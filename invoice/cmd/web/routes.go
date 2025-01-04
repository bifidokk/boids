package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	mux.Get("/", app.HomePage)
	mux.Get("/login", app.Login)
	mux.Post("/login", app.PostLogin)
	mux.Get("/register", app.Register)
	mux.Post("/register", app.PostRegister)
	mux.Post("/activate-account", app.ActivateAccount)
	mux.Get("/logout", app.Logout)

	return mux
}
