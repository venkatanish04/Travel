package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func NewDatabase(username, password, host, port, database string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open MySQL database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping MySQL database: %w", err)
	}
	return db, nil
}

func MySQLConfigFromEnv() MySQLConfig {
	return MySQLConfig{
		Username: envOrDefault("MYSQL_USER", "root"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     envOrDefault("MYSQL_HOST", "127.0.0.1"),
		Port:     envOrDefault("MYSQL_PORT", "3306"),
		Database: envOrDefault("MYSQL_DATABASE", "travelraft"),
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
