package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(writer http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(writer, r)
		return
	}

	writer.Write([]byte("Hello there!!"))
}

func showSnippet(writer http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.Error(writer, "No snippet found...", http.StatusNotFound)
		return
	}

	fmt.Fprintf(writer, "Displaying a specific snippet with id %d...", id)
}

func createSnippet(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	writer.Write([]byte("Creating a new snippet..."))
}
