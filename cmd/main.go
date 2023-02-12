package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/config"
	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/enforce"
)

func main() {
	if isTest := os.Getenv("IS_TEST"); isTest != "true" {
		if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
			log.Fatal().Err(err).Msg("failed to load .env file locally")
		}
	}
	cfg := config.Get()
	zerolog.SetGlobalLevel(cfg.LogLevel)
	if cfg.AppEnv != config.PRODUCTION {
		log.Debug().Msg("Starting local run")
		if err := enforce.Run(context.Background(), cfg); err != nil {
			log.Fatal().Err(err).Msg("failed to run locally")
		}
		log.Debug().Msg("Successfully finished running locally")
		return
	}
	lambda.Start(func(ctx context.Context) error {
		if err := enforce.Run(ctx, cfg); err != nil {
			return fmt.Errorf("failed to run: %w", err)
		}
		return nil
	})
}
