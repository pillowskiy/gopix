package domain

import "io"

// TEMP: Uncertain where to handle allowed extensions
// Currently handled in the infrastructure layer,
// but I lack confidence to rely on the infrastructure layer for this purpose. :|
var Unstable_AllowedContentTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/gif":  {},
	"image/webp": {},
	"image/avif": {},
	"image/bmp":  {},
	"image/tiff": {},

	"mp4":  {},
	"webm": {},
	"ogv":  {},
	"avi":  {},
	"mov":  {},
	"mkv":  {},
	"flv":  {},
	"wmv":  {},
}

type File struct {
	Reader io.ReadSeeker `json:"-"`
	Size   int64         `json:"-"`
}

type FileNode struct {
	File
	Name        string `json:"-"`
	ContentType string `json:"content_type,omitempty"`
}

func (f *FileNode) HasAllowedContentType() bool {
	_, ok := Unstable_AllowedContentTypes[f.ContentType]
	return ok
}
