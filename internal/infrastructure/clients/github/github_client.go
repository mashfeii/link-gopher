package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"
)

type HTTPRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	token      string
	endpoint   string
	httpClient HTTPRequestDoer
}

func NewClient(token, endpoint string, httpClient HTTPRequestDoer) *Client {
	return &Client{
		token:      token,
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

func (c *Client) GetUpdates(url string, since time.Time) ([]models.Event, error) {
	owner, repo, err := pkg.ValidateGithubURL(url)
	if err != nil {
		return nil, err
	}

	// First, get Pull Requests
	prURL := fmt.Sprintf("%s/repos/%s/%s/pulls?since=%s",
		c.endpoint, owner, repo, since.Format(time.RFC3339))

	prs, err := c.fetchEvents(prURL, "pull_request")
	if err != nil {
		return nil, err
	}

	// Then, get Issues
	issuesURL := fmt.Sprintf("%s/repos/%s/%s/issues?since=%s",
		c.endpoint, owner, repo, since.Format(time.RFC3339))

	issues, err := c.fetchEvents(issuesURL, "issue")
	if err != nil {
		return nil, err
	}

	return append(prs, issues...), nil
}

func (c *Client) fetchEvents(url string, eventType models.EventType) ([]models.Event, error) {
	const op = "github.Client.fetchEvents"

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %w", op, err)
	}

	req.Header.Add("Accept", "application/vnd.github+json")

	if c.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get %s: %w", op, eventType, err)
	}
	defer resp.Body.Close()

	var data []RepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("%s: failed to decode %s: %w", op, eventType, err)
	}

	events := make([]models.Event, len(data))
	for i, item := range data {
		events[i] = &RepoResponseWrapper{Item: item, Type: eventType}
	}

	return events, nil
}

// Blank implementation for GetQuestionTitle.
func (c *Client) GetQuestionTitle(_ string) (string, error) {
	return "", nil
}
