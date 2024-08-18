package storage

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type StorageHolder struct {
	Postgres *sqlx.DB
	Redis    *redis.Client
	S3       *s3.S3
	cfg      *config.Config
}

func NewStorageHolder(cfg *config.Config) *StorageHolder {
	return &StorageHolder{cfg: cfg}
}

func (s *StorageHolder) Setup() error {
	db, err := NewPostgres(&s.cfg.Postgres)
	if err != nil {
		return errors.Wrap(err, "storage.Setup.NewPostgres")
	}
	s.Postgres = db

	s.Redis = NewRedisClient(&s.cfg.Redis)

	s3, err := NewS3Storage(&s.cfg.S3)
	if err != nil {
		return errors.Wrap(err, "storage.Setup.NewS3Storage")
	}
	s.S3 = s3

	return nil
}

func (s *StorageHolder) Close() error {
	return s.Postgres.Close()
}
