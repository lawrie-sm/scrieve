package main

import (
	"errors"
	"fmt"
	tld "github.com/jpillora/go-tld"
	"html/template"
	"net/http"
	"os"
)

// validateURL validates a URL. Also escapes any special characters if needed.
// For speed, we don't check that it responds, just that it's valid
func validateURL(target string) (validated string, err error) {
	if len(target) > 2000 {
		err = errors.New("URL too long")
		return
	}

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

// getShortURL returns the full short URL path for a token
// by prepending the apps root URL
func getShortURL(token string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("ROOT_URL"), token)
}

// genHTML is a helper for template rendering
func genHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var paths []string
	for _, n := range filenames {
		paths = append(paths, fmt.Sprintf("web/views/%s.html", n))
	}
	template.Must(template.ParseFiles(paths...)).Execute(w, data)
}
