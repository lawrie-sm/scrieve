package main

import (
	"errors"
	"fmt"
	tld "github.com/jpillora/go-tld"
	"html/template"
	"net/http"
	"os"
)

// Validates a URL. Also escapes any special characters if needed.
// For speed, we don't check that it responds, just that it's valid
func validateURL(target string) (validated string, err error) {
	url, err := tld.Parse(target)
	if err != nil {
		return
	}

	// Add the scheme if its empty.
	// Re-parse to properly detect the other fields
	if url.Scheme == "" {
		url.Scheme = "http"
		url, err = tld.Parse(url.String())
		if err != nil {
			return
		}
	}

	// Validate the URL string by checking for a host and TLD
	if url.Host == "" || url.TLD == "" {
		return "", errors.New("Invalid URL")
	}

	validated = url.String()
	return
}

// Returns the full URL for a token, by prepending the app root URL
func getShortURL(token string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("ROOT_URL"), token)
}

// Helper to generate HTML
func genHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var paths []string
	for _, n := range filenames {
		paths = append(paths, fmt.Sprintf("web/views/%s.html", n))
	}
	template.Must(template.ParseFiles(paths...)).Execute(w, data)
}
