package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const (
	indexTemplate = "index.html"
)

type indexData struct{}

// Index is the front page of the web application
func (srv Web) Index(w http.ResponseWriter, r *http.Request) {
	data := indexData{}

	p := filepath.Join(defaultDocRoot, indexTemplate)
	t, err := template.ParseFiles(p)
	if err != nil {
		log.Printf("Error loading the application template: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
