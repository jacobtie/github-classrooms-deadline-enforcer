package config

import (
	"time"

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
	GitHub   struct {
		BaseURL        string        `envconfig:"GITHUB_BASE_URL" default:"http://localhost:3000/github"`
		AuthToken      string        `envconfig:"GITHUB_AUTH_TOKEN" default:"test-token"`
		Timeout        time.Duration `envconfig:"GITHUB_CONFIG" default:"500ms"`
		OrgName        string        `envconfig:"GITHUB_ORG_NAME" default:"test-org"`
		ConfigRepoName string        `envconfig:"GITHUB_CONFIG_REPO_NAME" default:"test-config-repo"`
	}
	Test struct {
		IS_TEST   bool   `envconfig:"IS_TEST" default:"false"`
		TEST_DATE string `envconfig:"TEST_DATE" default:""`
	}
}

func Get() *Config {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatal().Err(err).Msg("failed to parse env config")
	}
	return &c
}
