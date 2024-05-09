package main

import (
	"fmt"
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

	d := photodir.CrawlFilesystem(path)

	fmt.Printf("%v", d.Directories)
}
