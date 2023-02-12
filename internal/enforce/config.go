package enforce

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/config"
	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/github"
)

type enforcerConfig struct {
	Students    []*enforcerConfigStudent    `json:"students"`
	Assignments []*enforcerConfigAssignment `json:"assignments"`
}

type enforcerConfigStudent struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type enforcerConfigAssignment struct {
	Name       string                               `json:"name"`
	Deadline   string                               `json:"deadline"`
	Extensions []*enforcerConfigAssignmentExtension `json:"extensions,omitempty"`
}

type enforcerConfigAssignmentExtension struct {
	Name     string `json:"name"`
	Deadline string `json:"deadline"`
}

func getConfig(ctx context.Context, cfg *config.Config, gh *github.GitHubClient) (*enforcerConfig, error) {
	b, err := gh.GetFile(ctx, cfg.GitHub.ConfigRepoName, "config.json")
	if err != nil {
		return nil, fmt.Errorf("[enforcer.getConfig] failed to read config.json from repo: %w", err)
	}
	var enforcerCfg enforcerConfig
	if err := json.Unmarshal(b, &enforcerCfg); err != nil {
		return nil, fmt.Errorf("[enforcer.getConfig] failed to unmarshal config.json: %w", err)
	}
	return &enforcerCfg, nil
}
