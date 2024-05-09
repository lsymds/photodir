package main

import (
	"net/http"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/lsymds/photodir/internal/photodir"
)

func main() {
	path, err := filepath.Abs(".")
	if err != nil {
		log.Err(err).Msg("finding absolute path")
		return
	}

	// crawl the filesystem before booting the server
	d := photodir.CrawlFilesystem(path)

	// create the HTTP server
	h := photodir.NewWebServer(d)
	http.ListenAndServe("127.0.0.1:8994", h)
}
