package pkg

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func ShortenURL(rawURL string) string {
	re := regexp.MustCompile(`^(http[s]?://(www\.)?|ftp://(www\.)?|www\.)?([0-9A-Za-z.@:\%-_+~#=]+)(\.[a-zA-Z]{2,3})+(/.*)+$`)
	if !re.MatchString(rawURL) {
		return rawURL
	}

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
		user := Truncate(parts[0], 15)
		repo := Truncate(parts[1], 15)

		return fmt.Sprintf("git: %s:%s", user, repo)
	}

	if len(parts) == 1 {
		return "git: " + Truncate(parts[0], 25)
	}

	return "git: " + parsed.Host
}

func formatStackOverflowURL(parsed *url.URL) string {
	trimmedPath := strings.Trim(parsed.Path, "/")
	parts := strings.Split(trimmedPath, "/")

	if len(parts) >= 2 && parts[0] == "questions" {
		if len(parts) >= 3 {
			slug := strings.ReplaceAll(parts[2], "-", " ")
			return "stack: " + Truncate(slug, 22)
		}

		return "stack: " + parts[1]
	}

	return formatGenericURL(parsed)
}

func formatGenericURL(parsed *url.URL) string {
	hostParts := strings.Split(parsed.Host, ".")

	var domain string

	if len(hostParts) == 2 {
		domain = hostParts[0]
	} else if len(hostParts) > 2 {
		domain = hostParts[len(hostParts)-2]
	}

	if domain == "" {
		domain = parsed.Host
	}

	trimmedPath := strings.Trim(parsed.Path, "/")
	pathParts := strings.Split(trimmedPath, "/")
	path := ""

	if len(pathParts) > 0 {
		path = Truncate(pathParts[0], 15)
	}

	result := domain
	if path != "" {
		result += "/" + path
	}

	return Truncate(result, 30)
}

func Truncate(s string, maxChars int) string {
	if len(s) > maxChars {
		if maxChars <= 3 {
			return s[:maxChars]
		}

		return s[:maxChars-3] + "..."
	}

	return s
}
