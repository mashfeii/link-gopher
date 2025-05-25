package stackoverflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"
)

type HTTPRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	key        string
	endpoint   string
	httpClient HTTPRequestDoer
}

func NewClient(key, endpoint string, client HTTPRequestDoer) *Client {
	return &Client{
		key:        key,
		endpoint:   endpoint,
		httpClient: client,
	}
}

func (c *Client) GetUpdates(url string, since time.Time) ([]models.Event, error) {
	const op = "stackoverflow.Client.GetUpdates"

	// Validate the StackOverflow URL and extract the question ID
	questionID, err := pkg.ValidateStackOverflowURL(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errors.NewErrInvalidURL(url))
	}

	// First, get answers to the question
	answersURL := fmt.Sprintf("%s/questions/%d/answers?site=stackoverflow&filter=withbody&fromdate=%d", c.endpoint, questionID, since.Unix())

	answers, err := c.fetchEvents(answersURL, models.EventTypeAnswer)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to fetch answers: %w", op, err)
	}

	// Then, get comments on the question
	commentsURL := fmt.Sprintf("%s/questions/%d/comments?site=stackoverflow&filter=withbody&fromdate=%d",
		c.endpoint, questionID, since.Unix())

	comments, err := c.fetchEvents(commentsURL, models.EventTypeComment)
	if err != nil {
		return nil, err
	}

	return append(answers, comments...), nil
}

func (c *Client) fetchEvents(url string, format models.EventType) ([]models.Event, error) {
	const op = "stackoverflow.Client.fetchEvents"

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %w", op, err)
	}

	req.Header.Add("Accept", "application/json")

	if c.key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to do request: %w", op, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var events SQLQuestionResponse
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("%s: failed to decode response: %w", op, err)
	}

	typedEvents := make([]models.Event, 0, len(events.Items))
	for i, item := range events.Items {
		typedEvents[i] = &SQLQuestionItemWrapper{Item: item, Type: format}
	}

	return typedEvents, nil
}

func (c *Client) GetQuestionTitle(url string) (string, error) {
	const op = "stackoverflow.Client.GetQuestionTitle"

	questionID, err := pkg.ValidateStackOverflowURL(url)
	if err != nil {
		return "", err
	}

	reqURL := fmt.Sprintf("%s/questions/%d?site=stackoverflow&filter=withbody", c.endpoint, questionID)

	req, err := http.NewRequest(http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("%s: failed to create request: %w", op, err)
	}

	req.Header.Add("Accept", "application/json")

	if c.key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: failed to get question title: %w", op, err)
	}
	defer resp.Body.Close()

	var data SQLQuestionResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("%s: failed to decode question title: %w", op, err)
	}

	if len(data.Items) == 0 {
		return "", fmt.Errorf("%s: no items found for question ID %d", op, questionID)
	}

	return data.Items[0].GetTitle(), nil
}
