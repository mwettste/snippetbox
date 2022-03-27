package main

import (
	"net/http"

	"github.com/mwettste/snippetbox/pkg/fsmiddleware"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fsmiddleware.DisallowDirListing(fileServer)))

	return mux
}
