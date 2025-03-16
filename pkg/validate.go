package pkg

import (
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

	return matches[1], matches[2], nil
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

	return questionID, nil
}
