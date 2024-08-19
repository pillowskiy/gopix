package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
)

const (
	MaxFileSize = 300 * 1024 * 1024
)

var allowedImagesContentType = map[string]string{
	"image/png":  "png",
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
	"image/gif":  "gif",
	"image/webp": "webp",
	"image/avif": "avif",
	"video/mp4":  "mp4",
	"video/webm": "webm",
}

type FileNode struct {
	Data        []byte `json:"-"`
	Name        string `json:"-"`
	Size        int    `json:"-"`
	ContentType string `json:"-"`
}

func (n *FileNode) Prepare() error {
	if n.Size > MaxFileSize {
		return errors.New("file too big")
	}

	ext, err := n.Extension()
	if err != nil {
		return err
	}

	n.Name = fmt.Sprintf("%s.%s", n.randomName(), ext)
	return nil
}

func (n *FileNode) HasCorrectType() bool {
	_, allowed := allowedImagesContentType[n.ContentType]
	return allowed
}

func (n *FileNode) Extension() (string, error) {
	extension, exists := allowedImagesContentType[n.ContentType]
	if !exists {
		return "", errors.New("file extension not allowed")
	}
	return extension, nil
}

// TEMP: it's not a good solution
func (n *FileNode) randomName() string {
	rdn := make([]byte, 16)
	_, _ = rand.Read(rdn)
	return fmt.Sprintf("%x", rdn)
}
