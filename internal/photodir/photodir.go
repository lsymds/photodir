package photodir

import "image"

// ImageDirectory represents a file-system directory that contains photos. A tree of directories will be created for any
// nested directories that exist on the filesystem.
type ImageDirectory struct {
	Directories []ImageDirectory
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
	Thumbnail ImageThumbnail
}

// ImageThumbnail contains smaller representations of images.
type ImageThumbnail struct {
	Height int
	Width  int
	Image  *image.NRGBA
}
