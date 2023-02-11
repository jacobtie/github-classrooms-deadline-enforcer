package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Environment string

const (
	LOCAL      Environment = "local"
	PRODUCTION Environment = "production"
)

type Config struct {
	AppEnv   Environment   `envconfig:"APP_ENV" default:"local"`
	LogLevel zerolog.Level `envconfig:"LOG_LEVEL" default:"0"` // default to debug
}

var cfg *Config

func Get() *Config {
	if cfg != nil {
		return cfg
	}
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatal().Err(err).Msg("failed to parse env config")
	}
	cfg = &c
	return cfg
}
