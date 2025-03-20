package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

func ValidateGithubURL(url string) (owner, repo string, err error) {
	pattern := regexp.MustCompile(`^(?:https?://)?github\.com/([A-Za-z0-9-]+)/([A-Za-z0-9-]+)/?$`)
	matches := pattern.FindStringSubmatch(url)

	if len(matches) < 3 {
		return "", "", errors.NewErrInvalidURL(url)
	}

	owner, repo = matches[1], matches[2]
	err = validateGithubRepo(owner, repo)

	return
}

func ValidateStackOverflowURL(url string) (questionID int, err error) {
	pattern := regexp.MustCompile(`^(?:https?://)?stackoverflow\.com/questions/(\d+)/(.*)?`)
	matches := pattern.FindStringSubmatch(url)

	if len(matches) < 3 {
		return 0, errors.NewErrInvalidURL(url)
	}

	questionID, err = strconv.Atoi(matches[1])
	if err != nil {
		return 0, errors.NewErrInvalidURL(url)
	}

	err = validateStackQuestion(questionID)
	if err != nil {
		return 0, err
	}

	return questionID, nil
}

func validateGithubRepo(owner, repo string) error {
	endpoint := "https://api.github.com/repos"
	reqURL := fmt.Sprintf("%s/%s/%s/events", endpoint, owner, repo)

	req, _ := http.NewRequest("GET", reqURL, http.NoBody)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get events: %s", resp.Status)
	}

	return nil
}

func validateStackQuestion(questionID int) error {
	endpoint := "https://api.stackexchange.com/2.3/questions"
	reqURL := fmt.Sprintf("%s/%d?site=stackoverflow", endpoint, questionID)

	req, _ := http.NewRequest("GET", reqURL, http.NoBody)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()

	responseBody := struct {
		Items []struct{} `json:"items"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK || len(responseBody.Items) == 0 {
		return fmt.Errorf("failed to get question: %s", resp.Status)
	}

	return nil
}
