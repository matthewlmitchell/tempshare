package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError() writes the stack trace and error message to app.errorLog
// then forwards the client to a code 500 error page
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// Setting the call depth to 2 will print starting from the previous step of the trace,
	// since the last step in the trace will be the function that is printing the trace itself.
	app.errorLog.Output(2, trace)

	errorMessage := "The server encountered a problem and could not process your request."
	http.Error(w, errorMessage, http.StatusInternalServerError)
}

// clientError() responds to the client via an http responsewriter with an http error status code
func (app *application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

// runInBackground() accepts a function and runs it inside of a new goroutine
// while waiting to detect any panics. If a panic is detected in the goroutine,
// automatically recover and print the necessary trace to our app.errorLog
func (app *application) runInBackground(fn func()) {
	go func() {

		defer func() {
			if err := recover(); err != nil {
				app.errorLog.Printf("%s Trace: %s", err, debug.Stack())
			}
		}()

		fn()
	}()
}