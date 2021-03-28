package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var config Configuration

// Configuration define service config
type Configuration struct {
	fx.Out

	Env      string `yaml:"env" env:"ENVIRONMENT" default:"product"`
	Log      LogConfig
	Database DatabaseConfig
	HTTP     HTTPConfig
}

// New load App configuration
func New() (Configuration, error) {
	viper.AutomaticEnv()

	configPath := viper.GetString("CONFIG_PATH")
	if configPath == "" {
		configPath = "./deploy/config"
	}
	configName := viper.GetString("CONFIG_NAME")
	if configName == "" {
		configName = "app"
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Msgf("Error reading config file, err: %+v", err)
		return config, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, err: %+v", err)
		return config, err
	}

	return config, nil
}
