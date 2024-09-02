package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pillowskiy/gopix/internal/config"

	"github.com/aws/aws-sdk-go/service/s3"
	manager "github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	*s3.S3
	Uploader     *s3manager.Uploader
	PublicBucket string
}

func NewS3Storage(cfg *config.S3) (*S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretAccess, ""),
		S3ForcePathStyle: aws.Bool(cfg.ForcePathStyle),
	})

	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(
			cfg.UploadBufferSizeMB * 1024 * 1024,
		)
		u.PartSize = cfg.MultipartChunkSizeMB * 1024 * 1024
	})

	return &S3{
		S3:           s3.New(sess),
		Uploader:     uploader,
		PublicBucket: cfg.Bucket,
	}, nil
}
