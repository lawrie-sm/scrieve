package main

import (
	"html/template"
	"log"
	"net/http"
)

// serveIndex returns the standard index template
func (s *service) serveIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	genHTML(w, "", "index", "base")
	return
}

// serve400 renders the index page with an "invalid URL" error and 400 header
func (s *service) serve400(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	w.WriteHeader(http.StatusBadRequest)
	genHTML(w, "Invalid URL, please try again.", "index", "base")
	return
}

// serve404 returns the index page with a "not found" error and 404 header
func (s *service) serve404(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	w.WriteHeader(http.StatusBadRequest)
	genHTML(w, "That URL was not found.", "index", "base")
	return
}

// postURL processes the URL submission form, storing it in the database,
// getting a short link and rendering a page showing the short link
func (s *service) postURL(w http.ResponseWriter, r *http.Request) {
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

// serveRedirect handles shortened links. It looks them up in the database
// and issues an appropriate redirect
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

// handleRoot passes on all non root requests to be redirected
// and serves GET and POST requests on the root
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
		s.postURL(w, r)
		return
	}
	http.Error(w, "Method not allowed", 405)
	return
}

// setupRoutes creates a muxer and assigns handle functions
func (s *service) setupRoutes() {
	s.mux = http.NewServeMux()

	// File server for static assets
	fs := http.FileServer(http.Dir("web/static"))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Main handle func
	s.mux.HandleFunc("/", s.handleRoot)
}
