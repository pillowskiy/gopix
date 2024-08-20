package image

import (
	"fmt"
	"net/http"

	nanoid "github.com/matoous/go-nanoid/v2"
)

var imagesMimeTypeExt = map[string]string{
	"image/jpeg":    "jpg",
	"image/png":     "png",
	"image/gif":     "gif",
	"image/webp":    "webp",
	"image/avif":    "avif",
	"image/svg+xml": "svg",
	"image/x-icon":  "ico",
	"image/bmp":     "bmp",
	"image/tiff":    "tiff",

	"video/mp4":        "mp4",
	"video/webm":       "webm",
	"video/ogg":        "ogv",
	"video/x-msvideo":  "avi",
	"video/quicktime":  "mov",
	"video/x-matroska": "mkv",
	"video/x-flv":      "flv",
	"video/x-ms-wmv":   "wmv",
}

func GenerateUniqueFilename(ext string) string {
	str, err := nanoid.New(8)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s.%s", str, ext)
}

func GetMimeImageExt(data []byte) (string, error) {
	mime := DetectMimeFileType(data)
	ext, err := GetExtByMime(mime)
	if err != nil {
		return "", err
	}
	return ext, nil
}

func DetectMimeFileType(data []byte) string {
	return http.DetectContentType(data)
}

func GetExtByMime(mime string) (string, error) {
	ext, ok := imagesMimeTypeExt[mime]
	if !ok {
		return "", fmt.Errorf("unknown mime type: %s", mime)
	}
	return ext, nil
}
