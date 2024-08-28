package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message, use
	// the http.Error() function to send an Internal Server Error response to the
	// user, and then return from the handler so no subsequent code is executed.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Because the home handler is now a method against the application
		// struct it can access its fields, including the structured logger. We'll
		// use this to create a log entry at Error level containing the error
		// message, also including the request method and URI as attributes to
		// assist with debugging.
		app.serverError(w, r, err)
		return
	}

	// Then we use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Create some variables holding dummy data. We'll remove these later on
	// during the build.

	title := "0 mouse"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) downloadHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Clean("./ui/static/file.zip"))
}
