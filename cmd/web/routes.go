package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
	"github.com/schollz/httpfileserver"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := chi.NewRouter()
	mux.Get("/", app.home)
	mux.Get("/create", app.createTempShareForm)
	mux.Post("/create", app.createTempShare)


	fileServer := httpfileserver.New("/static/", "./ui/static/")
	mux.Get("/static/*", http.StripPrefix("/static", fileServer).(http.HandlerFunc))

	return standardMiddleware.Then(mux)
}