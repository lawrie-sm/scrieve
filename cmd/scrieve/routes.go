package main

import (
	"html/template"
	"log"
	"net/http"
)

// Render the index page
func (s *service) serveIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	genHTML(w, "", "index", "base")
	return
}

// Render the 400 (Bad Request) page - for invalid URLS
func (s *service) serve400(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	w.WriteHeader(http.StatusBadRequest)
	genHTML(w, "Invalid URL, please try again.", "index", "base")
	return
}

// Render the 404 (Not Found) page
func (s *service) serve404(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	w.WriteHeader(http.StatusBadRequest)
	genHTML(w, "That URL was not found.", "index", "base")
	return
}

// Render the shortened URL page
func (s *service) serveShortened(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)

	// Parse the form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Validate the URL
	target, err := validateURL(r.PostFormValue("target"))
	if err != nil {
		log.Println(err)
		s.serve400(w, r)
		return
	}

	// Create a new pair
	p, err := s.db.CreatePair(target)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Render the shortened page with the full URL and target
	type data struct {
		Short  template.URL
		Token  string
		Target template.URL
	}
	d := data{
		Short:  template.URL(getShortURL(p.Token)),
		Token:  p.Token,
		Target: template.URL(p.Target),
	}
	genHTML(w, d, "shortened", "base")

	return
}

// Lookup the short URL and issue a redirect
func (s *service) serveRedirect(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)

	// Grab the token from the end of the URL
	path := r.URL.Path
	token := path[1:len(path)]

	// Retrieve the target with the token
	target, err := s.db.GetTarget(token)
	if err != nil {
		log.Println(err)
		s.serve404(w, r)
		return
	}

	// Perform the redirect
	http.Redirect(w, r, target, 301)
	return
}

// Handles all interactions with the root - which is most of the app
func (s *service) handleRoot(w http.ResponseWriter, r *http.Request) {
	// Handle any redirects
	if len(r.URL.Path) > 1 {
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
