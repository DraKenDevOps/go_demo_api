package database

import (
	"fmt"
	"net/url"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go_demo_api/config"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn, err := parseDBURI(cfg.DBURI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB URI: %w", err)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

func parseDBURI(uri string) (string, error) {
	if !strings.HasPrefix(uri, "mariadb://") && !strings.HasPrefix(uri, "mysql://") {
		return uri, nil
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	username := parsed.User.Username()
	password, _ := parsed.User.Password()
	host := parsed.Hostname()
	port := parsed.Port()
	database := strings.TrimPrefix(parsed.Path, "/")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)

	return dsn, nil
}