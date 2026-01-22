package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"snippetbox.quasi.go/internal/models"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippetRepository.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	// Use the new render helper.
	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippetRepository.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Use the new render helper.
	app.render(w, r, http.StatusOK, "view.html", templateData{
		Snippet: snippet,
	})
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	snippet := models.Snippet{
		Title:   "O snail",
		Content: "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa",
		Expires: time.Now().AddDate(0, 0, 7),
	}

	id, err := app.snippetRepository.Insert(&snippet)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
