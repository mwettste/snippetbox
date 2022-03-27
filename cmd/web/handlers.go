package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(writer, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(writer, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
