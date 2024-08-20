package domain

// TEMP: Uncertain where to handle allowed extensions
// Currently handled in the infrastructure layer,
// but I lack confidence to rely on the infrastructure layer for this purpose. :|
var Unstable_AllowedFilesExt = map[string]struct{}{
	"png":  {},
	"jpg":  {},
	"jpeg": {},
	"gif":  {},
	"webp": {},
	"avif": {},
	"mp4":  {},
	"webm": {},
}

type FileNode struct {
	Data []byte `json:"-"`
	Name string `json:"-"`
	Size int64  `json:"-"`
}
