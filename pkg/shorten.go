package pkg

import (
	"fmt"
	"net/url"
	"strings"
)

func ShortenURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	switch parsed.Host {
	case "github.com":
		return formatGitHubURL(parsed)
	case "stackoverflow.com":
		return formatStackOverflowURL(parsed)
	default:
		return formatGenericURL(parsed)
	}
}

func formatGitHubURL(parsed *url.URL) string {
	trimmedPath := strings.Trim(parsed.Path, "/")
	parts := strings.Split(trimmedPath, "/")

	if len(parts) >= 2 {
		user := truncate(parts[0], 15)
		repo := truncate(parts[1], 15)

		return fmt.Sprintf("git: %s:%s", user, repo)
	}

	if len(parts) == 1 {
		return "git: " + truncate(parts[0], 25)
	}

	return "git: " + parsed.Host
}

func formatStackOverflowURL(parsed *url.URL) string {
	trimmedPath := strings.Trim(parsed.Path, "/")
	parts := strings.Split(trimmedPath, "/")

	if len(parts) >= 2 && parts[0] == "questions" {
		if len(parts) >= 3 {
			slug := strings.ReplaceAll(parts[2], "-", " ")
			return "stack: " + truncate(slug, 22)
		}

		return "stack: " + parts[1]
	}

	return formatGenericURL(parsed)
}

func formatGenericURL(parsed *url.URL) string {
	hostParts := strings.Split(parsed.Host, ".")

	var domain string

	for _, part := range hostParts {
		if part != "www" && part != "" {
			domain = part
			break
		}
	}

	if domain == "" {
		domain = parsed.Host
	}

	trimmedPath := strings.Trim(parsed.Path, "/")
	pathParts := strings.Split(trimmedPath, "/")
	path := ""

	if len(pathParts) > 0 {
		path = truncate(pathParts[0], 15)
	}

	result := domain
	if path != "" {
		result += "/" + path
	}

	return truncate(result, 30)
}

func truncate(s string, maxChars int) string {
	if len(s) > maxChars {
		if maxChars <= 3 {
			return s[:maxChars]
		}

		return s[:maxChars-3] + "..."
	}

	return s
}
