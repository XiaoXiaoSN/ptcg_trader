package database

import (
	"time"

	"ptcg_trader/internal/config"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabases init and return write and read DB objects
func InitDatabases(cfg config.DatabaseConfig) (*gorm.DB, error) {
	db, err := SetupDatabase(cfg)
	if err != nil {
		return nil, err
	}

	return db, err
}

// SetupDatabase ...
func SetupDatabase(database config.DatabaseConfig) (*gorm.DB, error) {
	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = time.Duration(180) * time.Second

	// get database connection string
	dsn, err := database.GetConnectionStr()
	if err != nil {
		return nil, err
	}

	var dialector gorm.Dialector
	switch database.Type {
	case config.MySQL:
		dialector = mysql.Open(dsn)
	case config.Postgres:
		dialector = postgres.Open(dsn)
	}

	// setting gorm logger
	logLevel := logger.Warn
	if database.Debug {
		logLevel = logger.Info
	}
	newLogger := NewLogger(logger.Config{
		SlowThreshold: time.Second, // Slow SQL threshold
		LogLevel:      logLevel,    // Log level
	})

	var conn *gorm.DB
	err = backoff.Retry(func() error {
		var err error

		gormConfig := gorm.Config{
			Logger: newLogger,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		}
		conn, err = gorm.Open(dialector, &gormConfig)
		if err != nil {
			log.Warn().Msgf("gorm failed to open database connection: %v", err)
			return err
		}

		sqlDB, err := conn.DB()
		if err != nil {
			log.Warn().Msgf("failed to get sql.db from gorm.DB instance: %v", err)
			return err
		}

		err = sqlDB.Ping()
		if err != nil {
			log.Warn().Msgf("ping db failed: %v", err)
			return err
		}

		return nil
	}, backOff)

	if err != nil {
		log.Error().Msgf("main: database connect err: %s", err.Error())
		return nil, err
	}
	log.Info().Msgf("database ping success")

	// setting sql connection
	{
		sqlDB, err := conn.DB()
		if err != nil {
			return nil, err
		}

		if database.MaxIdleConns != 0 {
			sqlDB.SetMaxIdleConns(database.MaxIdleConns)
		} else {
			sqlDB.SetMaxIdleConns(2)
		}

		if database.MaxOpenConns != 0 {
			sqlDB.SetMaxOpenConns(database.MaxOpenConns)
		} else {
			sqlDB.SetMaxOpenConns(5)
		}

		if database.MaxLifetimeSec != 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(database.MaxLifetimeSec) * time.Second)
		} else {
			sqlDB.SetConnMaxLifetime(14400 * time.Second)
		}
	}

	return conn, nil
}
