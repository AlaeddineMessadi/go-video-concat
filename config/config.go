package config

import (
	"errors"
	"fmt"
	"os"
)

type (
	Config struct {
		Database DatabaseConfig
		Storage  StorageConfig
		Redis    RedisConfig
	}
	DatabaseConfig struct {
		Host         string
		Username     string
		Password     string
		Name         string
		Port         string
	}
	StorageConfig struct {
		MinioUser     string
		MinioPassword string
		MinioHost     string
	}
	RedisConfig struct {
		Host     string
		Password string
		DB       uint
	}
)

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		},
		Storage: StorageConfig{
			MinioUser:     os.Getenv("MINIO_USER"),
			MinioPassword: os.Getenv("MINIO_PASSWORD"),
			MinioHost:     os.Getenv("MINIO_HOST"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0, // Optionally make this configurable
		},
	}

	// Validate required variables
	missing := []string{}
	if cfg.Database.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.Database.Username == "" {
		missing = append(missing, "DB_USERNAME")
	}
	if cfg.Database.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.Database.Name == "" {
		missing = append(missing, "DB_NAME")
	}
	if cfg.Database.Port == "" {
		missing = append(missing, "DB_PORT")
	}
	if cfg.Storage.MinioUser == "" {
		missing = append(missing, "MINIO_USER")
	}
	if cfg.Storage.MinioPassword == "" {
		missing = append(missing, "MINIO_PASSWORD")
	}
	if cfg.Storage.MinioHost == "" {
		missing = append(missing, "MINIO_HOST")
	}
	if cfg.Redis.Host == "" {
		missing = append(missing, "REDIS_HOST")
	}
	if len(missing) > 0 {
		return nil, errors.New(fmt.Sprintf("Missing required environment variables: %v", missing))
	}

	return cfg, nil
}
