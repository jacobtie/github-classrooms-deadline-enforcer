package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/config"
	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/enforcer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Get()
	zerolog.SetGlobalLevel(cfg.LogLevel)
	if cfg.AppEnv == config.LOCAL {
		log.Debug().Msg("Starting local run")
		if err := enforcer.Enforce(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("failed to run locally")
		}
		log.Debug().Msg("Successfully finished running locally")
		return
	}
	lambda.Start(enforcer.Enforce)
}
