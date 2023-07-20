package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

// handlerSPA implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type handlerSPA struct {
	staticPath string
	indexPath  string
}

func NewHandlerSPA(staticPath, indexPath string) handlerSPA {
	return handlerSPA{staticPath: staticPath, indexPath: indexPath}
}

func (receiver handlerSPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the absolute path to prevent directory traversal.
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Prepend the path with the path to the static directory.
	path = filepath.Join(receiver.staticPath, path)
	// Check if a file exists.
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(receiver.staticPath, receiver.indexPath))
		return
	} else if err != nil {
		// If we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(receiver.staticPath)).ServeHTTP(w, r)
}
