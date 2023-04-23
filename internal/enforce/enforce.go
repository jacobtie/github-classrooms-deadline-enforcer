package enforce

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/config"
	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/github"
	"github.com/rs/zerolog/log"
)

type enforcement struct {
	repoName string
	username string
}

func Run(ctx context.Context, cfg *config.Config) error {
	gh := github.NewGitHubClient(cfg)
	enforcerCfg, err := getConfig(ctx, cfg, gh)
	if err != nil {
		return fmt.Errorf("[enforcer.Enforce]: %w", err)
	}
	enforcements, err := resolveReposToEnforce(cfg, enforcerCfg)
	if err != nil {
		return fmt.Errorf("[enforcer.Enforce]: %w", err)
	}
	for _, enforcement := range enforcements {
		if err := gh.MakeUserReader(ctx, enforcement.repoName, enforcement.username); err != nil {
			errMsg := fmt.Sprintf("ERROR: failed to update user permissions to read for %s in repo %s: %s", enforcement.username, enforcement.repoName, err.Error())
			log.Error().Err(err).Msg(errMsg)
		}
		delay()
	}
	return nil
}

func delay() {
	// Create some randomness with our sleep
	jitterMS := rand.Intn(250) + 1
	sleepAmount := time.Duration(jitterMS+250) + 250
	time.Sleep(sleepAmount)
}
