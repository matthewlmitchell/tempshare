package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	defaultMiddleware := alice.New(app.recoverPanic, app.logRequest)

	mux := chi.NewRouter()
	mux.Get("/", app.home)

	return defaultMiddleware.Then(mux)
}