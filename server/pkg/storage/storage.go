package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/config"
	"github.com/redis/go-redis/v9"
)

type StorageHolder struct {
	Postgres *sqlx.DB
	Redis    *redis.Client
	cfg      *config.Config
}

func NewStorageHolder(cfg *config.Config) *StorageHolder {
	return &StorageHolder{cfg: cfg}
}

func (s *StorageHolder) Setup() error {
	db, err := NewPostgres(&s.cfg.Postgres)
	if err != nil {
		return err
	}
	s.Postgres = db

	s.Redis = NewRedisClient(&s.cfg.Redis)

	return nil
}

func (s *StorageHolder) Close() error {
	return s.Postgres.Close()
}
