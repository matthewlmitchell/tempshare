package main

import (
	"fmt"
	"net/http"

	"github.com/matthewlmitchell/tempshare/pkg/forms"
	"github.com/matthewlmitchell/tempshare/pkg/models"
	"github.com/matthewlmitchell/tempshare/pkg/recaptcha"
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
	form.Required("text", "expires", "viewlimit", "g-recaptcha-response")
	form.MinLength("text", 2)
	form.MaxLength("text", 1024)
	form.PermittedValues("expires", "1", "3", "7")
	form.PermittedValues("viewlimit", "1", "3", "10")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	success, err := recaptcha.VerifyRecaptcha(app.serverConfig.env, app.httpsClient, r, form.Get("g-recaptcha-response"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	if !success {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	tempShare, err := app.tempShare.New(form.Get("text"), form.Get("expires"), form.Get("viewlimit"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", fmt.Sprintf("Your TempShare link: %s", fmt.Sprintf("https://placeholder.com/view?token=%s", tempShare.PlainText)))

	// Refresh the page so the message flash will become visible
	//http.Redirect(w, r, "/create", http.StatusSeeOther)

	app.render(w, r, "create.page.tmpl", &templateData{Form: form})

}

func (app *application) viewTempShareForm(w http.ResponseWriter, r *http.Request) {

	formData := forms.New(r.URL.Query())

	app.render(w, r, "view.page.tmpl", &templateData{Form: formData})
}

func (app *application) viewTempShare(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("token", "g-recaptcha-response")
	form.MaxLength("token", 52)
	form.MinLength("token", 52)

	if !form.Valid() {
		form.Errors.Add("generic", "Invalid token")
		app.render(w, r, "view.page.tmpl", &templateData{Form: form})
		return
	}

	success, err := recaptcha.VerifyRecaptcha(app.serverConfig.env, app.httpsClient, r, form.Get("g-recaptcha-response"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	if !success {
		app.session.Put(r, "flash", "An error occurred.\nPlease complete the captcha again.")
		app.render(w, r, "view.page.tmpl", &templateData{Form: form})
		return
	}

	token := &models.TempShare{
		PlainText: form.Get("token"),
	}

	tempShareData, err := app.tempShare.Get(token.PlainText)
	if err == models.ErrNoRecord {
		form.Errors.Add("generic", "Invalid token")
		app.render(w, r, "view.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// TODO: Delete tempShare from database when views >= viewlimit

	app.session.Put(r, "flash", fmt.Sprintf("This link has %d uses remaining.", tempShareData.ViewLimit-tempShareData.Views-1))

	app.render(w, r, "home.page.tmpl", &templateData{TempShare: tempShareData})
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "about.page.tmpl", nil)
}
