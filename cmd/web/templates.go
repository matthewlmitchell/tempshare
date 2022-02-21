package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/matthewlmitchell/tempshare/pkg/forms"
	"github.com/matthewlmitchell/tempshare/pkg/models"
)

type templateData struct {
	CurrentYear int
	CSRFToken   string
	Flash       string
	TempShare   *models.TempShare
	Form        *forms.Form
}

// FormattedDate accepts a time.Time and returns the time adjusted to UTC.
// If the time given is zero, an empty string is returned
func FormattedDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("Jan 02 2006 at 15:04")
}

// This template.FuncMap{} allows us to call Golang functions
// inside our template files, e.g.: {{formattedDate .VariableName}}
var functions = template.FuncMap{
	"formattedDate": FormattedDate,
}

// initTemplateCache accepts a directory and returns a map that points
// file names to parsed template.Template objects.
func initTemplateCache(dir string) (map[string]*template.Template, error) {

	// Create a new map to hold templates as a cache
	cache := map[string]*template.Template{}

	// Search the directory and return all filepaths ending in ".page.tmpl"
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// For every page template found in the directory:
	for _, page := range pages {
		fileName := filepath.Base(page)

		// .New(name).Funcs(functions).ParseFiles(page) will be used
		// when we have some functions that need to be executed inside of the template
		// functions is a template.FuncMap{}

		// Create a new HTML template with the filename above
		templateParsed, err := template.New(fileName).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Locate any layout templates, parse them, and add them to the set of templates
		templateParsed, err = templateParsed.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// The same as above, but for partial templates
		templateParsed, err = templateParsed.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the parsed template into the cache map
		cache[fileName] = templateParsed
	}

	return cache, nil
}
