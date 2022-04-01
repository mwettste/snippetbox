package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
	app.render(writer, r, "createsnippet.page.tmpl", nil)
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

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "The title must not be empty"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "Title is too long (max. 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "Content must not be empty"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		app.render(writer, r, "createsnippet.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
