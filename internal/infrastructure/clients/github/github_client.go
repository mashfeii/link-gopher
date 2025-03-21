package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"
)

type Client struct {
	token      string
	endpoint   string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		endpoint:   "https://api.github.com/repos",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetEvent(url string) (models.Event, error) {
	owner, repo, err := pkg.ValidateGithubURL(url)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s/%s/%s/events", c.endpoint, owner, repo)

	req, err := http.NewRequest("GET", reqURL, http.NoBody)
	req.Header.Add("Accept", "application/vnd.github+json")

	if c.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get events: %s", resp.Status)
	}

	var events []RepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &events[0], nil
}
