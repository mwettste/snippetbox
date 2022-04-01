package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mwettste/snippetbox/pkg/forms"
	"github.com/mwettste/snippetbox/pkg/models"
)

func (app *application) home(writer http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) createSnippetForm(writer http.ResponseWriter, r *http.Request) {
	app.render(writer, r, "createsnippet.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) showSnippet(writer http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}

		return
	}

	app.render(writer, r, "snippet.page.tmpl", &templateData{
		Snippet: snippet,
	})
}

func (app *application) createSnippet(writer http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(writer, r, "createsnippet.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
