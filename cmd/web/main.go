package main

import (
	"log"
	"net/http"

	"github.com/mwettste/snippetbox/pkg/fsmiddleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fsmiddleware.DisallowDirListing(fileServer)))

	log.Println("Starting server on :4000")
	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}
