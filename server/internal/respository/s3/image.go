package s3

import (
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/storage"
)

const multipartUploadThreshold = 5 * 1024 * 1024

type imageStorage struct {
	s3     *storage.S3
	bucket string
}

func NewImageStorage(s3 *storage.S3, bucket string) *imageStorage {
	return &imageStorage{s3: s3, bucket: bucket}
}

func (s *imageStorage) Put(ctx context.Context, file *domain.FileNode) error {
	reader := bytes.NewReader(file.Data)
	defer reader.Reset([]byte(""))

	if file.Size < multipartUploadThreshold {
		_, err := s.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(file.Name),
			Body:   reader,
		})
		return err
	}

	log.Printf(
		"File size is greater than %vMB (is: %vMV) -> Multipart upload",
		multipartUploadThreshold/1024/1024,
		file.Size/1024/1024,
	)

	_, err := s.s3.Uploader.UploadWithContext(ctx, &manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(file.Name),
		Body:   reader,
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