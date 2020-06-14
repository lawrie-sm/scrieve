package main

import (
	"html/template"
	"log"
	"net/http"
)

// Render the index page
func (s *service) handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)

	tmpl, err := template.ParseFiles("web/views/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = tmpl.Execute(w, "")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}

	return
}

func (s *service) setupRoutes() {
	s.mux = http.NewServeMux()

	// File server for static assets
	fs := http.FileServer(http.Dir("web/static"))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Main page
	s.mux.HandleFunc("/", s.handleIndex)
}
