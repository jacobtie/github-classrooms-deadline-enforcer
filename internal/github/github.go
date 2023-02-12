package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/jacobtie/github-classrooms-deadline-enforcer/internal/config"
)

type GitHubClient struct {
	httpClient *http.Client
	baseURL    string
	authToken  string
	orgName    string
}

type requestParams struct {
	requestMethod      string
	requestURL         string
	body               io.Reader
	expectedStatusCode int
	extraHeaders       map[string]string
	discardResponse    bool
}

func NewGitHubClient(cfg *config.Config) *GitHubClient {
	c := &GitHubClient{
		httpClient: &http.Client{
			Timeout: cfg.GitHub.Timeout,
		},
		baseURL:   cfg.GitHub.BaseURL,
		authToken: cfg.GitHub.AuthToken,
		orgName:   cfg.GitHub.OrgName,
	}
	return c
}

func (gh *GitHubClient) makeRequest(ctx context.Context, params *requestParams) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, params.requestMethod, params.requestURL, params.body)
	if err != nil {
		return nil, fmt.Errorf("[github.makeRequest] failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gh.authToken))
	for name, val := range params.extraHeaders {
		req.Header.Add(name, val)
	}
	resp, err := gh.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[github.makeRequest] failed to do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != params.expectedStatusCode {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("[github.makeRequest] failed with status %s and could not ready body with error: %w", resp.Status, err)
		}
		return nil, fmt.Errorf("[github.makeRequest] failed with status %s and body %s", resp.Status, string(b))
	}
	if params.discardResponse {
		// Read the response to EOF for keep-alive
		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			return nil, fmt.Errorf("[github.makeRequest] failed to read response which should not happen")
		}
		return nil, nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[github.makeRequest] failed to read body: %w", err)
	}
	return b, nil
}
