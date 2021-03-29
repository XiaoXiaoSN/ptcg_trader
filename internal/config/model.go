package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// LogConfig ...
type LogConfig struct {
	Level       string `mapstructure:"level" env:"LOG_LEVEL"`
	Environment string `mapstructure:"environment" env:"LOG_ENVIRONMENT"`
	Format      string `mapstructure:"format" env:"LOG_FORMAT"`
	AppID       string `mapstructure:"app_id" env:"LOG_APP_ID"`
}

// DatabaseType type define
type DatabaseType string

const (
	// MySQL DatabaseType
	MySQL DatabaseType = "mysql"
	// Postgres DatabaseType
	Postgres DatabaseType = "postgres"
)

// DatabaseConfig for db config
type DatabaseConfig struct {
	Debug          bool         `mapstructure:"debug" env:"DB_DEBUG"`
	Type           DatabaseType `mapstructure:"type" env:"DB_TYPE"`
	Host           string       `mapstructure:"host" env:"DB_HOST"`
	Port           int          `mapstructure:"port" env:"DB_PORT"`
	Username       string       `mapstructure:"username" env:"DB_USERNAME"`
	Password       string       `mapstructure:"password" env:"DB_PASSWORD"`
	DBName         string       `mapstructure:"db_name" env:"DB_NAME"`
	MaxIdleConns   int          `mapstructure:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns   int          `mapstructure:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	MaxLifetimeSec int          `mapstructure:"max_lifetime_sec" env:"DB_MAX_LIFETIME_SEC"`
	// pg should setting this value. It will restrict access to db schema.
	SearchPath string `mapstructure:"search_path" env:"DB_SEARCH_PATH"`
	// pg ssl mode
	SSLEnable bool `mapstructure:"ssl_enable" env:"DB_SSL_ENABLE"`
}

// GetConnectionStr ...
func (database *DatabaseConfig) GetConnectionStr() (string, error) {
	var connectionString string
	switch database.Type {
	case MySQL:
		connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&multiStatements=true", database.Username, database.Password, database.Host+":"+strconv.Itoa(database.Port), database.DBName)
	case Postgres:
		connectionString = fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s`, database.Username, database.Password, database.Host, database.Port, database.DBName)

		if database.SSLEnable {
			connectionString += " sslmode=require"
		} else {
			connectionString += " sslmode=disable"
		}

		if strings.TrimSpace(database.SearchPath) != "" {
			connectionString = fmt.Sprintf("%s search_path=%s", connectionString, database.SearchPath)
		}
	default:
		return "", errors.New("Not support driver")
	}
	return connectionString, nil
}

// HTTPConfig setting http config
type HTTPConfig struct {
	Mode    string `mapstructure:"mode" env:"HTTP_MODE"` // ex: debug, product
	Address string `mapstructure:"address" env:"HTTP_ADDRESS"`
}
