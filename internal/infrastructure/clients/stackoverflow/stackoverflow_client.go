package stackoverflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"
)

type Client struct {
	key        string
	endpoint   string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		key:        apiKey,
		endpoint:   "https://api.stackexchange.com/2.3/questions",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetEvent(url string) (models.Event, error) {
	questionID, err := pkg.ValidateStackOverflowURL(url)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s/%d?site=stackoverflow", c.endpoint, questionID)
	if c.key != "" {
		reqURL += fmt.Sprintf("&key=%s", c.key)
	}

	req, err := http.NewRequest("GET", reqURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StackOverflow API error: %s", resp.Status)
	}

	var questionData SOQuestionResponse
	if err := json.NewDecoder(resp.Body).Decode(&questionData); err != nil {
		return nil, err
	}

	if len(questionData.Items) == 0 {
		return nil, fmt.Errorf("question not found")
	}

	return &questionData, nil
}
