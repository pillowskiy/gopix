package rest

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pkg/errors"
)

const fileHeaderSize int64 = 512

var (
	errNotAllowedImageExt = NewBadRequestError("Prohibited image extension")
	errInvalidImageData   = NewBadRequestError("Invalid image data")
)

func ReadEchoImage(c echo.Context, field string) (*domain.FileNode, error) {
	image, err := c.FormFile(field)
	if err != nil {
		return nil, errInvalidImageData
	}

	if err := checkImageContentType(image); err != nil {
		return nil, err
	}

	file, err := image.Open()
	if err != nil {
		return nil, errors.Wrap(err, "ReadEchoImage.Open")
	}
	defer file.Close()

	binImage := bytes.NewBuffer(nil)
	if _, err := io.Copy(binImage, file); err != nil {
		return nil, errors.Wrap(err, "ReadEchoImage.ReadFrom")
	}

	node := &domain.FileNode{
		Data:        binImage.Bytes(),
		Name:        image.Filename,
		Size:        int(image.Size),
		ContentType: image.Header["Content-Type"][0],
	}

	return node, nil
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
