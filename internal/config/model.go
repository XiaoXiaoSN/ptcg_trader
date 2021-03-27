package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// LogConfig ...
type LogConfig struct {
	Level       string `yaml:"level" env:"level"`
	Environment string `yaml:"environment" env:"environment"`
	Format      string `yaml:"format" env:"format"`
	AppID       string `yaml:"app_id" env:"app_id"`
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
	Debug          bool
	Type           DatabaseType
	Host           string
	Port           int
	Username       string
	Password       string
	DBName         string
	MaxIdleConns   int
	MaxOpenConns   int
	MaxLifetimeSec int
	// pg should setting this value. It will restrict access to db schema.
	SearchPath string `yaml:"search_path" mapstructure:"search_path"`
	// pg ssl mode
	SSLEnable bool `yaml:"ssl_enable" mapstructure:"ssl_enable"`
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
	Mode    string `json:"mode"` // ex: debug, product
	Address string `json:"address"`
}
