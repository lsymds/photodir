package photodir

import (
	"errors"
	"image"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
	"github.com/rs/zerolog/log"
)

// CrawlFilesystem traverses the filesystem under the given path identifying any images that should be shown.
func CrawlFilesystem(path string) *Directory {
	root := &Directory{
		Path: path,
		Name: path,
	}

	// find all images
	log.Info().Msg("crawling directories for images")
	crawlDirectory(path, root)

	// generate their thumbnails concurrently
	log.Info().Msg("generating thumbnails")
	wg := &sync.WaitGroup{}
	generateThumbnails(root, wg)
	wg.Wait()

	return root
}

// crawlDirectory recursively iterates over a given directory finding child directories (which are also recursively
// iterated through) and image files
func crawlDirectory(path string, parent *Directory) {
	log.Debug().Str("path", path).Msg("crawling directory")

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		log.Err(err).Str("path", path).Msg("reading directory")
		return
	}

	for _, e := range dirEntries {
		entryPath := filepath.Join(path, e.Name())

		// if the entry is a directory, then recursively crawl it and its children
		if e.IsDir() {
			d := &Directory{
				Name: e.Name(),
				Path: entryPath,
			}

			crawlDirectory(d.Path, d)

			// only pollute the directory tree with this directory if there is an image file within it
			if dirHasImages(d) {
				parent.Directories = append(parent.Directories, *d)
			}
		} else {
			// else, add any image files
			o, err := os.Open(entryPath)
			defer func() {
				o.Close()
			}()

			if err != nil {
				log.Err(err).Str("path", entryPath).Msg("opening path for reading")
				continue
			}

			// preferable to Decode as it doesn't result in the entire image being read into memory
			img, _, err := image.DecodeConfig(o)
			if err != nil {
				if errors.Is(err, image.ErrFormat) {
					log.Debug().Str("path", entryPath).Msg("file not a valid image")
				} else {
					log.Err(err).Str("path", entryPath).Msg("decoding image")
				}

				continue
			}

			f := ImageFile{
				Name:   e.Name(),
				Path:   entryPath,
				Height: img.Height,
				Width:  img.Width,
			}

			parent.ImageFiles = append(parent.ImageFiles, f)
		}
	}
}

// dirHasImages identifies whether there are any images present in the directory tree
func dirHasImages(d *Directory) bool {
	if len(d.ImageFiles) > 0 {
		return true
	}

	for _, cd := range d.Directories {
		if dirHasImages(&cd) {
			return true
		}
	}

	return false
}

// generateThumbnails concurrently generates thumbnail images for all image files in the directory tree
func generateThumbnails(d *Directory, wg *sync.WaitGroup) {
	// dispatch child goroutines to generate thumbnails for all directories
	for _, cd := range d.Directories {
		generateThumbnails(&cd, wg)
	}

	// generate the actual thumbnails in coroutines too
	for _, f := range d.ImageFiles {
		wg.Add(1)
		go generateThumbnail(&f, wg)
	}
}

// generateThumbnail generates a thumbnail image for a single image file
func generateThumbnail(f *ImageFile, wg *sync.WaitGroup) {
	defer wg.Done()

	t1 := time.Now()

	log.Info().Str("path", f.Path).Msg("generating thumbnail")

	fl, err := os.Open(f.Path)
	if err != nil {
		log.Err(err).Str("path", f.Path).Msg("opening file")
		return
	}

	img, _, err := image.Decode(fl)
	if err != nil {
		log.Err(err).Str("path", f.Path).Msg("decoding image")
		return
	}

	rimg := imaging.Resize(img, 360, 0, imaging.Lanczos)

	f.Thumbnail = Thumbnail{
		Width:  rimg.Bounds().Dx(),
		Height: rimg.Bounds().Dy(),
		Image:  rimg,
	}

	t2 := time.Now()

	log.
		Debug().
		Str("path", f.Path).
		Int("width", f.Thumbnail.Width).
		Int("height", f.Thumbnail.Height).
		Durs("duration", []time.Duration{t2.Sub(t1)}).
		Msg("generated thumbnail")
}
