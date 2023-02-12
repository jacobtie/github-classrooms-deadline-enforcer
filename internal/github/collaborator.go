package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (gh *GitHubClient) MakeUserReader(ctx context.Context, repoName, username string) error {
	requestURL := fmt.Sprintf("%s/repos/%s/%s/collaborators/%s", gh.baseURL, gh.orgName, repoName, username)
	if _, err := gh.makeRequest(ctx, &requestParams{
		requestMethod:      http.MethodPut,
		requestURL:         requestURL,
		body:               bytes.NewReader(json.RawMessage(`{"permission":"pull"}`)),
		expectedStatusCode: http.StatusNoContent,
		discardResponse:    true,
	}); err != nil {
		return fmt.Errorf("[github.removeCollaborator]: %w", err)
	}
	return nil
}
