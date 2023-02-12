package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (gh *GitHubClient) MakeUserReader(ctx context.Context, repoName, username string) error {
	if err := gh.removeCollaborator(ctx, repoName, username); err != nil {
		return fmt.Errorf("[github.MakeUserReader]: %w", err)
	}
	if err := gh.addCollaboratorAsReader(ctx, repoName, username); err != nil {
		return fmt.Errorf("[github.MakeUserReader]: %w", err)
	}
	return nil
}

func (gh *GitHubClient) removeCollaborator(ctx context.Context, repoName, username string) error {
	requestURL := fmt.Sprintf("%s/repos/%s/%s/collaborators/%s", gh.baseURL, gh.orgName, repoName, username)
	if _, err := gh.makeRequest(ctx, &requestParams{
		requestMethod:      http.MethodDelete,
		requestURL:         requestURL,
		expectedStatusCode: http.StatusNoContent,
		discardResponse:    true,
	}); err != nil {
		return fmt.Errorf("[github.addCollaboratorAsReader]: %w", err)
	}
	return nil
}

func (gh *GitHubClient) addCollaboratorAsReader(ctx context.Context, repoName, username string) error {
	requestURL := fmt.Sprintf("%s/repos/%s/%s/collaborators/%s", gh.baseURL, gh.orgName, repoName, username)
	if _, err := gh.makeRequest(ctx, &requestParams{
		requestMethod:      http.MethodPut,
		requestURL:         requestURL,
		body:               bytes.NewReader(json.RawMessage(`{"permission":"pull"}`)),
		expectedStatusCode: http.StatusCreated,
		discardResponse:    true,
	}); err != nil {
		return fmt.Errorf("[github.removeCollaborator]: %w", err)
	}
	return nil
}
