package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pillowskiy/gopix/internal/config"

	"github.com/aws/aws-sdk-go/service/s3"
)

func NewS3Storage(cfg *config.S3) (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretAccess, ""),
		S3ForcePathStyle: aws.Bool(cfg.ForcePathStyle),
	})

	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}
