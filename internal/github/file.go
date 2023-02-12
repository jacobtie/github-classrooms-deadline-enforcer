package github

import (
	"context"
	"fmt"
	"net/http"
)

func (gh *GitHubClient) GetFile(ctx context.Context, repoName, filename string) ([]byte, error) {
	requestURL := fmt.Sprintf("%s/repos/%s/%s/contents/%s", gh.baseURL, gh.orgName, repoName, filename)
	b, err := gh.makeRequest(ctx, &requestParams{
		requestMethod:      http.MethodGet,
		requestURL:         requestURL,
		expectedStatusCode: http.StatusOK,
		extraHeaders:       map[string]string{"Accept": "application/vnd.github.raw"}, // read file contents instead of metadata
	})
	if err != nil {
		return nil, fmt.Errorf("[github.GetFile]: %w", err)
	}
	return b, nil
}
