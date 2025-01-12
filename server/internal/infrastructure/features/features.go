package features

import (
	"context"
	"fmt"
	"io"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/image"
)

type basicFeatureExtractor struct{}

func NewBasicFeatureExtractor() *basicFeatureExtractor {
	return &basicFeatureExtractor{}
}

func (e *basicFeatureExtractor) MakeFileNode(ctx context.Context, file *domain.File) (*domain.FileNode, error) {
	contentType, err := image.DetectMimeFileType(file.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to detect mime file type: %w", err)
	}

	ext, err := image.GetExtByMime(contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to get extension by mime: %w", err)
	}

	return &domain.FileNode{
		File:        *file,
		Name:        image.GenerateUniqueFilename(ext),
		ContentType: contentType,
	}, nil
}

func (e *basicFeatureExtractor) Features(ctx context.Context, fileNode *domain.FileNode) (imgProps *domain.ImageProperties, err error) {
	info, err := image.GetImageInfo(fileNode.Reader.(io.ReaderAt))
	if err != nil {
		err = fmt.Errorf("failed to get image info: %w", err)
		return
	}

	imgProps = &domain.ImageProperties{
		Width:  info.Width,
		Height: info.Height,
		Ext:    info.Format,
	}

	return
}
