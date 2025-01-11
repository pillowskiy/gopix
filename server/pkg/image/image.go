package image

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pillowskiy/imagesize"
)

var imagesMimeTypeExt = map[string]string{
	"image/jpeg":       "jpg",
	"image/png":        "png",
	"image/gif":        "gif",
	"image/webp":       "webp",
	"image/avif":       "avif",
	"image/svg+xml":    "svg",
	"image/x-icon":     "ico",
	"image/bmp":        "bmp",
	"image/tiff":       "tiff",
	"video/mp4":        "mp4",
	"video/webm":       "webm",
	"video/ogg":        "ogv",
	"video/x-msvideo":  "avi",
	"video/quicktime":  "mov",
	"video/x-matroska": "mkv",
	"video/x-flv":      "flv",
	"video/x-ms-wmv":   "wmv",
}

type ImageInfo struct {
	Width  int
	Height int
	Format string
}

const (
	fallbackMime         = "application/octet-stream"
	uniqueFilenameLength = 8
)

func GenerateUniqueFilename(ext string) string {
	str, err := nanoid.New(uniqueFilenameLength)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s.%s", str, ext)
}

func GetExtByMime(mime string) (string, error) {
	ext, ok := imagesMimeTypeExt[mime]
	if !ok {
		return "", errors.New("unsupported mime provided")
	}
	return ext, nil
}

func DetectMimeFileType(reader io.ReadSeeker) (mime string, err error) {
	if _, err = reader.Seek(0, io.SeekStart); err != nil {
		err = fmt.Errorf("failed to seek start: %w", err)
		return
	}

	defer func() {
		if _, err = reader.Seek(0, io.SeekStart); err != nil {
			err = fmt.Errorf("failed to seek start: %w", err)
		}
	}()

	var data [512]byte
	if _, err := io.ReadFull(reader, data[:]); err != nil {
		return fallbackMime, nil
	}

	mime = http.DetectContentType(data[:])
	return
}

func GetImageInfo(reader io.ReaderAt) (*ImageInfo, error) {
	info, err := imagesize.ExtractInfo(reader)
	if err != nil {
		return nil, err
	}

	return &ImageInfo{
		Width:  info.Width,
		Height: info.Height,
		Format: info.Format,
	}, nil
}
