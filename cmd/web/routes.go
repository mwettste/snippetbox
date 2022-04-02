package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"github.com/mwettste/snippetbox/pkg/fsmiddleware"
)

func (app *application) routes() http.Handler {
	stdMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fsmiddleware.DisallowDirListing(fileServer)))

	return stdMiddleware.Then(mux)
}
