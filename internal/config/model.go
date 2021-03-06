package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// LogConfig ...
type LogConfig struct {
	Level       string `yaml:"level" env:"LOG_LEVEL"`
	Environment string `yaml:"environment" env:"LOG_ENVIRONMENT"`
	Format      string `yaml:"format" env:"LOG_FORMAT"`
	NoColor     bool   `yaml:"no_color" env:"LOG_NO_COLOR"`
	AppID       string `yaml:"app_id" env:"LOG_APP_ID"`
}

// DatabaseType type define
type DatabaseType string

const (
	// MySQL DatabaseType
	MySQL DatabaseType = "mysql"
	// Postgres DatabaseType
	Postgres DatabaseType = "postgres"
)

// DatabaseConfig for db connection config
type DatabaseConfig struct {
	Debug          bool         `yaml:"debug" env:"DB_DEBUG"`
	Type           DatabaseType `yaml:"type" env:"DB_TYPE"`
	Host           string       `yaml:"host" env:"DB_HOST"`
	Port           int          `yaml:"port" env:"DB_PORT"`
	Username       string       `yaml:"username" env:"DB_USERNAME"`
	Password       string       `yaml:"password" env:"DB_PASSWORD"`
	DBName         string       `yaml:"db_name" env:"DB_NAME"`
	MaxIdleConns   int          `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns   int          `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	MaxLifetimeSec int          `yaml:"max_lifetime_sec" env:"DB_MAX_LIFETIME_SEC"`
	// pg should setting this value. It will restrict access to db schema.
	SearchPath string `yaml:"search_path" env:"DB_SEARCH_PATH"`
	// pg ssl mode
	SSLEnable bool `yaml:"ssl_enable" env:"DB_SSL_ENABLE"`
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

// RedisConfig setting redis connection config
type RedisConfig struct {
	Addresses  []string `yaml:"addresses" env:"REDIS_ADDRESSES"`
	Password   string   `yaml:"password" env:"REDIS_PASSWORD"`
	MaxRetries int      `yaml:"max_retries" env:"REDIS_MAX_RETRIES"`
	PoolSize   int      `yaml:"pool_size" env:"REDIS_POOL_SIZE"`
	DB         int      `yaml:"db" env:"REDIS_DB"`
}

// HTTPConfig setting http config
type HTTPConfig struct {
	Debug   bool   `yaml:"debug" env:"HTTP_DEBUG"`
	Address string `yaml:"address" env:"HTTP_ADDRESS"`
}

// TraderStrategy ...
type TraderStrategy string

// TraderStrategy enum
var (
	TraderStrategyUnknown               TraderStrategy = ""
	TraderStrategyDatabaseRowLock       TraderStrategy = "database_row_lock"
	TraderStrategyRedisLock             TraderStrategy = "redis_lock"
	TraderStrategyAsyncInMemoryMatching TraderStrategy = "async_in_memory_matching"
)

// TraderConfig ...
type TraderConfig struct {
	Strategy TraderStrategy `yaml:"strategy" env:"TRADER_STRATEGY"`
}

// KafkaConfig for kafka client
type KafkaConfig struct {
	Debug          bool     `yaml:"debug" env:"KAFKA_DEBUG"`
	BrokerList     []string `yaml:"address" env:"KAFKA_ADDRESS"`
	ClientID       string   `yaml:"client_id" env:"KAFKA_CLIENT_ID"`
	GroupName      string   `yaml:"group_name" env:"KAFKA_GROUP_NAME"`
	OffsetsInitial int64    `yaml:"offsets_initial" env:"KAFKA_OFFSETS_INITIAL"`
}

// StanConfig for stan client
type StanConfig struct {
	Debug       bool   `yaml:"debug" env:"STAN_DEBUG"`
	Address     string `yaml:"address" env:"STAN_ADDRESS"`
	ClusterID   string `yaml:"cluster_id" env:"STAN_CLUSTER_ID"`
	ClientID    string `yaml:"client_id" env:"STAN_CLIENT_ID"`
	DurableName string `yaml:"durable_name" env:"STAN_DURABLE_NAME"`
}
