package pkg_test

import (
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/pkg"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"short", 10, "short"},
		{"thisisaverylongstring", 10, "thisisa..."},
		{"exactlength", 11, "exactlength"},
		{"small", 3, "sma"},
		{"tiny", 2, "ti"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := pkg.Truncate(tt.input, tt.max)
			if got != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.expected)
			}
		})
	}
}

func TestShortenURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Invalid URL",
			input:    "invalid-url",
			expected: "invalid-url",
		},
		{
			name:     "Generic URL with no path",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "Generic URL with www and path",
			input:    "https://www.example.com/path/to/page",
			expected: "example/path",
		},
		{
			name:     "Generic URL without www",
			input:    "https://mydomain.org/about",
			expected: "mydomain/about",
		},
		{
			name:     "GitHub URL with two segments",
			input:    "https://github.com/user/repository",
			expected: "git: user:repository",
		},
		{
			name:     "GitHub URL with one segment",
			input:    "https://github.com/username",
			expected: "git: username",
		},
		{
			name:     "GitHub URL with extra path segments",
			input:    "https://github.com/user/repo/extra/path",
			expected: "git: user:repo",
		},
		{
			name:     "GitHub URL with empty path",
			input:    "https://github.com/",
			expected: "git: ",
		},
		{
			name:     "StackOverflow URL with slug",
			input:    "https://stackoverflow.com/questions/12345678/how-to-code",
			expected: "stack: how to code",
		},
		{
			name:     "StackOverflow URL without slug",
			input:    "https://stackoverflow.com/questions/87654321/",
			expected: "stack: 87654321",
		},
		{
			name:     "StackOverflow non-question URL",
			input:    "https://stackoverflow.com/users/12345/user",
			expected: "stackoverflow/users",
		},
		{
			name:     "Generic URL with long path",
			input:    "https://sub.domain.com/averyveryverylongpath",
			expected: "domain/averyveryver...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pkg.ShortenURL(tt.input)
			if got != tt.expected {
				t.Errorf("ShortenURL(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
