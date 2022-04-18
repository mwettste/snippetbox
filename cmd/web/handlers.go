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

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(writer, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(writer http.ResponseWriter, r *http.Request) {
	app.render(writer, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(writer http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(writer, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use!")
			app.render(writer, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(writer, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(writer, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(writer http.ResponseWriter, r *http.Request) {
	app.render(writer, r, "login.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) loginUser(writer http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or password is incorrect!")
			app.render(writer, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(writer, err)
		}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)

	redirectPath := app.session.PopString(r, "redirectPathAfterLogin")
	if redirectPath != "" {
		http.Redirect(writer, r, redirectPath, http.StatusSeeOther)
		return
	}

	http.Redirect(writer, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(writer http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been signed out successfully!")
	http.Redirect(writer, r, "/", http.StatusSeeOther)
}

func (app *application) userProfile(writer http.ResponseWriter, r *http.Request) {
	userId := app.session.Get(r, "authenticatedUserID").(int)
	user, err := app.users.Get(userId)
	if err != nil {
		app.serverError(writer, err)
	}

	app.render(writer, r, "profile.page.tmpl", &templateData{AuthenticatedUser: user})
}

func (app *application) changePasswordForm(writer http.ResponseWriter, r *http.Request) {
	app.render(writer, r, "changepassword.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) changePassword(writer http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("currentPassword", "newPassword", "newPasswordConfirmation")
	form.MinLength("newPassword", 10)
	if form.Get("newPassword") != form.Get("newPasswordConfirmation") {
		form.Errors.Add("newPasswordConfirmation", "Passwords do not match")
	}

	if !form.Valid() {
		app.render(writer, r, "changepassword.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	userID := app.session.GetInt(r, "authenticatedUserID")
	err = app.users.ChangePassword(userID, form.Get("currentPassword"), form.Get("newPassword"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("currentPassword", "Current password is incorrect")
			app.render(writer, r, "changepassword.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(writer, err)
		}

		return
	}

	app.session.Put(r, "flash", "Your password has been updated!")
	http.Redirect(writer, r, "/user/profile", http.StatusSeeOther)
}

func (app *application) about(writer http.ResponseWriter, r *http.Request) {
	app.render(writer, r, "about.page.tmpl", &templateData{})
}

func ping(writer http.ResponseWriter, r *http.Request) {
	writer.Write([]byte("OK"))
}
