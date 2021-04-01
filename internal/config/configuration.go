package config

import (
	"os"
	"path/filepath"
	"ptcg_trader/internal/errors"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

var config Configuration

// Configuration define service config
type Configuration struct {
	fx.Out

	Env      string `yaml:"env" env:"ENVIRONMENT" default:"product"`
	Log      LogConfig
	Database DatabaseConfig
	Redis    RedisConfig
	HTTP     HTTPConfig
}

// New load App configuration
func New() (Configuration, error) {
	var fileName, rootDirPath string
	if fileName = os.Getenv("CONFIG_NAME"); fileName == "" {
		fileName = "app.yaml"
	}
	if rootDirPath = os.Getenv("CONFIG_PATH"); rootDirPath == "" {
		rootDirPath = "./deploy/config"
	}

	configPath := filepath.Join(rootDirPath, fileName)
	if _, err := os.Stat(configPath); err != nil {
		log.Error().Msgf("Error reading config file, err: %+v", err)
		return config, errors.Wrap(errors.ErrInternalError, err.Error())
	}

	// Enable debug mode or set env `CONFIGOR_DEBUG_MODE` to true when running your application
	err := configor.New(&configor.Config{Debug: false}).Load(&config, configPath)
	if err != nil {
		log.Error().Msgf("Error set config file, err: %+v", err)
		return config, errors.Wrap(errors.ErrInternalError, err.Error())
	}

	return config, nil
}
