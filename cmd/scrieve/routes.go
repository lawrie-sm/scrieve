package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Helper to generate HTML
func genHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var paths []string
	for _, n := range filenames {
		paths = append(paths, fmt.Sprintf("web/views/%s.html", n))
	}
	template.Must(template.ParseFiles(paths...)).Execute(w, data)
}

// Render the index page
func (s *service) serveIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	genHTML(w, "", "index", "base")
	return
}

// Render the shortened URL page
func (s *service) serveShortened(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)

	// Parse the form and store the target URL
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// TODO: Add HTTP/HTTPS to this
	target := r.PostFormValue("full-url")

	// Create a new pair in the DB
	p, err := s.db.CreatePair(target)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Create the full URL for display and generate the template
	fullURL := fmt.Sprintf("%s/%s", os.Getenv("ROOT_URL"), p.Token)
	type data struct {
		FullURL string
		Target  string
	}
	d := data{fullURL, target}
	genHTML(w, d, "shortened", "base")

	return
}

// Lookup the short URL and issue a redirect
func (s *service) serveRedirect(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)

	// Grab the token from the end of the URL
	path := r.URL.String()
	token := path[1:len(path)]

	// Retrieve the target with the token
	target, err := s.db.GetTarget(token)
	if err != nil {
		log.Println(err)
		// TODO: Nicer 404 page
		http.Redirect(w, r, "/", 404)
		return
	}

	// Perform the redirect
	// We assume HTTP will be converted to HTTP in most cases
	http.Redirect(w, r, target, 302)
	return
}

// Handles all interactions with the root - which is most of the app
func (s *service) handleRoot(w http.ResponseWriter, r *http.Request) {
	// Handle any redirects
	if r.URL.String() != "/" {
		s.serveRedirect(w, r)
		return
	}
	// Handle the methods on the root
	if r.Method == http.MethodGet {
		s.serveIndex(w, r)
		return
	}
	if r.Method == http.MethodPost {
		s.serveShortened(w, r)
		return
	}
	http.Error(w, "Method not allowed", 405)
	return
}

func (s *service) setupRoutes() {
	s.mux = http.NewServeMux()

	// File server for static assets
	fs := http.FileServer(http.Dir("web/static"))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Main handle func
	s.mux.HandleFunc("/", s.handleRoot)
}
