package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

// HandlerSPA implements the http.Handler interface to serve a Single Page Application (SPA).
// It responds to HTTP requests using the given staticPath as the root path for the web directory
// and indexPath as the index file within that web directory.
type HandlerSPA struct {
	staticPath string // The path to the root of the web directory
	indexPath  string // The path to the index file within the web directory
}

// NewHandlerSPA creates a new instance of HandlerSPA with the provided staticPath and indexPath.
func NewHandlerSPA(staticPath, indexPath string) HandlerSPA {
	return HandlerSPA{staticPath: staticPath, indexPath: indexPath}
}

// ServeHTTP serves the HTTP request by serving static files and the SPA.
func (receiver HandlerSPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the absolute path to prevent directory traversal.
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// If we failed to get the absolute path, respond with a 400 bad request and stop.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepend the path with the path to the web directory.
	path = filepath.Join(receiver.staticPath, path)

	// Check if a file exists.
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, "Daaamn! It's not there.", http.StatusNotFound)
		return
	} else if err != nil {
		// If we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Otherwise, use http.FileServer to serve the web directory.
	http.FileServer(http.Dir(receiver.staticPath)).ServeHTTP(w, r)
}
