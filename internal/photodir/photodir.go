package photodir

import "image"

// Directory represents a file-system directory that contains photos. A tree of directories will be created for any
// nested directories that exist on the filesystem.
type Directory struct {
	Directories []Directory
	ImageFiles  []ImageFile
	Path        string
	Name        string
}

// ImageFile contains details about images contained within the filesystem.
type ImageFile struct {
	Path      string
	Name      string
	Height    int
	Width     int
	Thumbnail Thumbnail
}

// Thumbnail contains smaller representations of images.
type Thumbnail struct {
	Height int
	Width  int
	Image  *image.NRGBA
}
