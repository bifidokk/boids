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
	mux.Get("/activate", app.ActivateAccount)
	mux.Get("/logout", app.Logout)

	mux.Mount("/members", app.authRouter())

	return mux
}

func (app *Config) authRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.Auth)

	mux.Get("/plans", app.PlansList)
	mux.Get("/subscribe", app.Subscribe)

	return mux
}
