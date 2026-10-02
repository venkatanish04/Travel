package config

import "os"

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		User:     envOrDefault("MYSQL_USER", "root"),
		Password: envOrDefault("MYSQL_PASSWORD", "root"),
		Host:     envOrDefault("MYSQL_HOST", "127.0.0.1"),
		Port:     envOrDefault("MYSQL_PORT", "3306"),
		Name:     envOrDefault("MYSQL_DATABASE", "travelraft"),
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
