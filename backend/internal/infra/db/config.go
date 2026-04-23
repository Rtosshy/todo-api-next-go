package database

import "backend/pkg"

type Config struct {
	Host     string
	Database string
	Port     string
	Driver   string
	User     string
	Password string
	SSLMode  string
}

func NewConfigPostgres() *Config {
	return &Config{
		Host:     pkg.GetEnvDefault("DB_HOST", "localhost"),
		Database: pkg.GetEnvDefault("DB_NAME", "api_database"),
		Port:     pkg.GetEnvDefault("DB_PORT", "5432"),
		Driver:   pkg.GetEnvDefault("DB_DRIVER", "postgres"),
		User:     pkg.GetEnvDefault("DB_USER", "app"),
		Password: pkg.GetEnvDefault("DB_PASSWORD", "password"),
		SSLMode:  pkg.GetEnvDefault("DB_SSL_MODE", "disable"),
	}
}

func NewConfigSQLite() *Config {
	return &Config{
		Database: pkg.GetEnvDefault("DB_NAME", "api_database.sqlite"),
	}
}
