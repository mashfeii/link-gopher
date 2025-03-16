package pkg_test

import (
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/pkg"
)

func TestValidateGithubURL(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedOwner string
		expectedRepo  string
		expectError   bool
	}{
		{
			name:          "Valid GitHub URL",
			input:         "https://github.com/user/repo",
			expectedOwner: "user",
			expectedRepo:  "repo",
			expectError:   false,
		},
		{
			name:          "Valid GitHub URL with hyphen",
			input:         "https://github.com/some-user/some-repo",
			expectedOwner: "some-user",
			expectedRepo:  "some-repo",
			expectError:   false,
		},
		{
			name:        "Invalid GitHub URL: missing repo",
			input:       "https://github.com/user",
			expectError: true,
		},
		{
			name:        "Invalid GitHub URL: wrong domain",
			input:       "https://notgithub.com/user/repo",
			expectError: true,
		},
		{
			name:        "Invalid GitHub URL: invalid characters",
			input:       "https://github.com/user_/repo_",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := pkg.ValidateGithubURL(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if owner != tt.expectedOwner || repo != tt.expectedRepo {
					t.Errorf("For input %q, expected (%q, %q) but got (%q, %q)",
						tt.input, tt.expectedOwner, tt.expectedRepo, owner, repo)
				}
			}
		})
	}
}

func TestValidateStackOverflowURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedID  int
		expectError bool
	}{
		{
			name:        "Valid StackOverflow URL",
			input:       "https://stackoverflow.com/questions/12345678/how-to-code",
			expectedID:  12345678,
			expectError: false,
		},
		{
			name:        "Valid StackOverflow URL with trailing slash",
			input:       "https://stackoverflow.com/questions/87654321/",
			expectedID:  87654321,
			expectError: false,
		},
		{
			name:        "Invalid StackOverflow URL: non-numeric ID",
			input:       "https://stackoverflow.com/questions/abc/how-to-code",
			expectError: true,
		},
		{
			name:        "Invalid StackOverflow URL: missing question segment",
			input:       "https://stackoverflow.com/questions/",
			expectError: true,
		},
		{
			name:        "Invalid StackOverflow URL: missing slash after ID",
			input:       "https://stackoverflow.com/questions/12345678",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := pkg.ValidateStackOverflowURL(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if id != tt.expectedID {
					t.Errorf("For input %q, expected questionID %d but got %d", tt.input, tt.expectedID, id)
				}
			}
		})
	}
}
