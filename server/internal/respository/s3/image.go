package s3

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pillowskiy/gopix/internal/domain"
)

type imageStorage struct {
	s3     *s3.S3
	bucket string
}

func NewImageStorage(s3 *s3.S3, bucket string) *imageStorage {
	return &imageStorage{s3: s3, bucket: bucket}
}

func (s *imageStorage) Put(ctx context.Context, file *domain.FileNode) error {
	_, err := s.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(file.Name),
		Body:   bytes.NewReader(file.Data),
	})
	return err
}

func (s *imageStorage) Delete(ctx context.Context, path string) error {
	_, err := s.s3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	return err
}
