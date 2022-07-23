package db

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Conn   *gorm.DB
	Logger *zap.Logger
}

func Connect(cfg Config, logger *zap.Logger) (*DB, error) {
	dsn := "host=" + cfg.Host + "user=" + cfg.User + "password=" + cfg.Password + "dbname=" + cfg.DBName + "port=" + fmt.Sprintf("%v", cfg.Port) + "sslmode=" + cfg.SSLMode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &DB{}, err
	}
	return &DB{Conn: db, Logger: logger}, err
}
