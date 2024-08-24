package storage

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/config"
)

const (
	maxOpenConns    = 10
	connMaxLifeTime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

const (
	SSLDisable = "disable"
	SSLEnable  = "enable"
)

func NewPostgres(cfg *config.Postgres) (*sqlx.DB, error) {
	sslStr := getStringSSL(cfg.SSL)
	dataSourceName := fmt.Sprintf(
		"host=%s port=%v user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Database, sslStr, cfg.Password,
	)

	db, err := sqlx.Connect(cfg.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifeTime * time.Second)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getStringSSL(ssl bool) string {
	if !ssl {
		return SSLDisable
	}
	return SSLEnable
}
