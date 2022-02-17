package main

import (
	"fmt"
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
	form.Required("text", "expires", "viewlimit")
	form.MinLength("text", 2)
	form.MaxLength("text", 1024)
	form.PermittedValues("expires", "1", "3", "7")
	form.PermittedValues("viewlimit", "1", "3", "100000")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// TODO: Require captcha verification or account registration to submit
	tempShare, err := app.tempShare.New(form.Get("text"), form.Get("expires"), form.Get("viewlimit"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// TODO: Generate a unique URL, insert the data into the database, and
	// return the URL to the user
	app.session.Put(r, "flash", fmt.Sprintf("Your TempShare link: %s", fmt.Sprintf("https://placeholder.com/view?token=%s", tempShare.PlainText)))

	// Refresh the page so the message flash will become visible
	http.Redirect(w, r, "/create", http.StatusSeeOther)

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "about.page.tmpl", nil)
}
