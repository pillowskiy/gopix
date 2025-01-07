package image

import (
	"fmt"
	"io"
	"net/http"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pillowskiy/imagesize"
)

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
