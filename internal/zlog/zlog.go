package zlog

import (
	"fmt"
	"os"

	"ptcg_trader/internal/config"
	"ptcg_trader/internal/errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLog init the log system setting
func InitLog(logConfig config.LogConfig) (err error) {
	zerolog.DisableSampling(true)

	hostname, err := os.Hostname()
	if err != nil {
		return errors.WithStack(err)
	}

	// set the log level
	logLevel, err := zerolog.ParseLevel(logConfig.Level)
	if err != nil {
		log.Warn().Msgf("%+v", errors.WithStack(err))
	}
	// with default log level: info
	if logLevel == zerolog.NoLevel {
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.
		New(os.Stdout).
		Level(logLevel).
		With().
		Timestamp().
		Fields(map[string]interface{}{
			"app_id":      logConfig.AppID,
			"hostname":    hostname,
			"environment": logConfig.Environment,
		}).
		Logger()

	// set the log format
	switch logConfig.Format {
	case "console":
		output := zerolog.ConsoleWriter{
			Out: os.Stdout,
			FormatTimestamp: func(v interface{}) string {
				return fmt.Sprintf("%s", v)
			},
		}
		logger = logger.Output(output)

	default:
		// zerolog default value (json format)
	}

	// 設定預設 logger
	log.Logger = logger

	return nil
}
