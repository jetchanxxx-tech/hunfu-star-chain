package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/huifu/star-chain/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL(cfg config.DatabaseConfig) (*sql.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("mysql dsn is empty")
	}
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("mysql open: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("mysql ping: %w", err)
	}
	return db, nil
}
