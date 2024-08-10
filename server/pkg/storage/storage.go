package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/config"
)

type StorageHolder struct {
	Postgres *sqlx.DB
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
	return nil
}

func (s *StorageHolder) Close() error {
	return s.Postgres.Close()
}
