package photodir

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

// recoveryMiddleware recovers any panics, logging and redirecting the consumer to the oops page.
func recoveryMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context())

		defer func() {
			if err := recover(); err != nil {
				l.Error().Any("panic", err).Msg("recovered from panic")
				http.Redirect(w, r, "/oops", http.StatusSeeOther)
			}
		}()

		h.ServeHTTP(w, r)
	})
}

// loggingMiddleware creates and assigns to the request's context a logger with properties extracted from the request
func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := log.With().
			Str("url", r.URL.String()).
			Str("method", r.Method)

		r = r.WithContext(l.Logger().WithContext(r.Context()))

		h.ServeHTTP(w, r)
	})
}

// ServeHTTP implements the [http.Handler] interface.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recoveryMiddleware(
		loggingMiddleware(
			s.router,
		),
	).ServeHTTP(w, r)
}

// handleGetIndex serves the root page for the application.
func (s *server) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	pageIndex().Render(r.Context(), w)
}
