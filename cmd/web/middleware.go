package main

import (
	"fmt"
	"net/http"
)

// recoverPanic defines a deferred function that will run in the event of a panic 
// on the application's main thread, which will attempt to automatically recover from the
// panic, log the error to app.errorLog, and close the connection to the client as code 500
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// This deferred function will run when a panic occurs
		defer func() {
			// If a panic is detected in our main thread: close the connection,
			// log the error to our app.errorLog logger, then forward the client
			// to a generic code 500 error page
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// logRequest logs where the request came from, the protocol of the request and its HTTP method,
// and what URL the request was for to our app.infoLog, then forwards the client's request
// to the next http.Handler
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Log where the request came from, the protocol of the request and its HTTP method,
		// and what URL the request was for
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}