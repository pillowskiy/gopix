package rest

import (
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	errNotAllowedImageExt = NewBadRequestError("Prohibited image extension")
	errInvalidImageData   = NewBadRequestError("Invalid image data")
)

func ReadEchoImage(c echo.Context, field string) (*multipart.FileHeader, error) {
	image, err := c.FormFile(field)

	if err != nil {
		return nil, errInvalidImageData
	}

	if err := checkImageContentType(image); err != nil {
		return nil, err
	}

	return image, nil
}

func checkImageContentType(file *multipart.FileHeader) error {
	contentType, err := determineFileContentType(file.Header)
	if err != nil {
		return errInvalidImageData
	}

	isImage := strings.HasPrefix(contentType, "image")
	isVideo := strings.HasPrefix(contentType, "video")
	if !isImage && !isVideo {
		return errNotAllowedImageExt
	}

	return nil
}

func determineFileContentType(fileHeader textproto.MIMEHeader) (string, error) {
	contentTypes := fileHeader["Content-Type"]
	if len(contentTypes) < 1 {
		return "", errNotAllowedImageExt
	}

	return contentTypes[0], nil
}
