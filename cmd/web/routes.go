package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
	"github.com/schollz/httpfileserver"
)

func (app *application) routes() http.Handler {

	// TODO: Add rate limiting to our dynamicMiddlware, before enabling http session management
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noCSRF)

	mux := chi.NewRouter()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home).(http.HandlerFunc))
	mux.Get("/create", dynamicMiddleware.ThenFunc(app.createTempShareForm).(http.HandlerFunc))
	mux.Post("/create", dynamicMiddleware.ThenFunc(app.createTempShare).(http.HandlerFunc))

	mux.Get("/view", dynamicMiddleware.ThenFunc(app.viewTempShareForm).(http.HandlerFunc))
	mux.Post("/view", dynamicMiddleware.ThenFunc(app.viewTempShare).(http.HandlerFunc))

	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about).(http.HandlerFunc))

	// TODO: Add rate limiting to the http file server
	fileServer := httpfileserver.New("/static/", "./ui/static/")
	mux.Get("/static/*", http.StripPrefix("/static", fileServer).(http.HandlerFunc))

	return standardMiddleware.Then(mux)
}
