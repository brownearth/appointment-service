package config

import (
	"fmt"
	"os"
	"strconv"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Test        Environment = "test"
)

type StorageType string

const (
	Memory   StorageType = "memory"
	Postgres StorageType = "postgres"
	SqlLite3 StorageType = "sqlite3"
)

type Config struct {
	Environment    Environment
	LogLevel       string
	LogSource      bool
	LogFormat      string
	StorageType    StorageType
	SqlLite3DbFile string
	Port           string
	DB             DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func Load() *Config {
	return &Config{
		Environment:    Environment(envOrDefault("APP_ENV", "development")),
		Port:           envOrDefault("APP_PORT", "8080"),
		LogLevel:       envOrDefault("LOG_LEVEL", "info"),
		LogSource:      envAsBool("LOG_SOURCE", true),
		LogFormat:      envOrDefault("LOG_FORMAT", "text"),
		StorageType:    StorageType(envOrDefault("STORAGE_TYPE", "memory")),
		SqlLite3DbFile: envOrDefault("DB_FILE", ""),
		DB: DBConfig{
			Host:     envOrDefault("DB_HOST", ""),
			Port:     envOrDefault("DB_PORT", ""),
			Name:     envOrDefault("DB_NAME", ""),
			User:     envOrDefault("DB_USER", ""),
			Password: envOrDefault("DB_PASSWORD", ""),
		},
	}
}

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func envAsBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"=============================================================\n"+
			"Configuration\n"+
			"-------------------------------------------------------------\n"+
			"Config{\n"+
			"  Environment: %s\n"+
			"  LogLevel: %s\n"+
			"  LogSource: %t\n"+
			"  LogFormat: %s\n"+
			"  StorageType: %s\n"+
			"  SqlLite3DbFile: %s\n"+
			"  Port: %s\n"+
			"  DB: {\n"+
			"    Host: %s\n"+
			"    Port: %s\n"+
			"    Name: %s\n"+
			"    User: %s\n"+
			"    Password: ***\n"+ // Hide password
			"  }\n"+
			"}\n"+
			"=============================================================",
		c.Environment,
		c.LogLevel,
		c.LogSource,
		c.LogFormat,
		c.StorageType,
		c.SqlLite3DbFile,
		c.Port,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
		c.DB.User,
	)
}
