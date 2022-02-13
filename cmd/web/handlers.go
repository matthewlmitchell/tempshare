package main

import (
	"net/http"

	"github.com/matthewlmitchell/tempshare/pkg/forms"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "home.page.tmpl", nil)
}

func (app *application) createTempShareForm(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createTempShare(w http.ResponseWriter, r *http.Request) {

	// Parse the HTTP POST request for data to populate r.PostForm and r.Form .
	// If any errors return, tell the client their request was bad.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("text", "expires")
	form.MinLength("text", 2)
	form.MaxLength("text", 1024)
	form.PermittedValues("expires", "0", "1", "7")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// TODO: Require captcha verification to submit

	// TODO: Generate a unique URL, insert the data into the database, and
	// return the URL to the user

	// tempShareURL := ""

	// Add session handling
	// app.session.Put(r, "flash", fmt.Sprintf("Your TempShare link: %s", tempShareURL))

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "about.page.tmpl", nil)
}
