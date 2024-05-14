package photodir

import (
	"net/http"

	"github.com/lsymds/go-utils/pkg/http/middleware"
	"github.com/rs/zerolog"
)

// server is a wrapper around the HTTP server
type server struct {
	router *http.ServeMux
	dir    *ImageDirectory
}

// NewWebServer creates a web server that serves the website of the application.
func NewWebServer(dir *ImageDirectory) http.Handler {
	s := &server{
		router: &http.ServeMux{},
		dir:    dir,
	}

	s.router.HandleFunc("GET /", s.handleGetIndex)

	return s
}

// ServeHTTP implements the [http.Handler] interface.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.Recovery(
		middleware.Logging(
			s.router,
			func(c *zerolog.Context) {},
		),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/oops", http.StatusSeeOther)
		}),
	).ServeHTTP(w, r)
}

// handleGetIndex serves the root page for the application.
func (s *server) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	pageIndex().Render(r.Context(), w)
}
