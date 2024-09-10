package domain

// TEMP: Uncertain where to handle allowed extensions
// Currently handled in the infrastructure layer,
// but I lack confidence to rely on the infrastructure layer for this purpose. :|
var Unstable_AllowedContentTypes = map[string]struct{}{
	"image/jpeg":    {},
	"image/png":     {},
	"image/gif":     {},
	"image/webp":    {},
	"image/avif":    {},
	"image/svg+xml": {},
	"image/x-icon":  {},
	"image/bmp":     {},
	"image/tiff":    {},

	"mp4":  {},
	"webm": {},
	"ogv":  {},
	"avi":  {},
	"mov":  {},
	"mkv":  {},
	"flv":  {},
	"wmv":  {},
}

type FileNode struct {
	Data        []byte `json:"-"`
	Name        string `json:"-"`
	Size        int64  `json:"-"`
	ContentType string `json:"content_type,omitempty"`
}

func (f *FileNode) HasAllowedContentType() bool {
	_, ok := Unstable_AllowedContentTypes[f.ContentType]
	return ok
}
