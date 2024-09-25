/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"testing"
)

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		url            string
		expectedOwner  string
		expectedRepo   string
		expectedBranch string
		expectedPath   string
		expectError    bool
	}{
		// {
		// 	url:            "https://github.com/stuttgart-things/stuttgart-things.git@main:kaeffken/apps/flux/app-defaults.yaml",
		// 	expectedOwner:  "stuttgart-things",
		// 	expectedRepo:   "stuttgart-things",
		// 	expectedBranch: "main",
		// 	expectedPath:   "kaeffken/apps/flux/app-defaults.yaml",
		// 	expectError:    false,
		// },
		// {
		// 	// Test case with missing branch
		// 	url:            "https://github.com/stuttgart-things/stuttgart-things.git@:kaeffken/apps/flux/app-defaults.yaml",
		// 	expectedOwner:  "",
		// 	expectedRepo:   "",
		// 	expectedBranch: "",
		// 	expectedPath:   "",
		// 	expectError:    true,
		// },
		// {
		// Test case with invalid format (no '@')
		// 	url:         "https://github.com/stuttgart-things/stuttgart-things.git:kaeffken/apps/flux/app-defaults.yaml",
		// 	expectError: true,
		// },
		// {
		// 	// Test case with invalid format (no ':')
		// 	url:         "https://github.com/stuttgart-things/stuttgart-things.git@mainkaeffken/apps/flux/app-defaults.yaml",
		// 	expectError: true,
		// },
		// {
		// 	// Test case with no owner/repo format
		// 	url:         "https://github.com/stuttgart-things/.git@main:kaeffken/apps/flux/app-defaults.yaml",
		// 	expectError: true,
		// },
	}

	for _, test := range tests {
		owner, repo, branch, path, err := ParseGitHubURL(test.url)
		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for URL: %s, but got none", test.url)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for URL: %s: %v", test.url, err)
			}
			if owner != test.expectedOwner {
				t.Errorf("Expected owner: %s, but got: %s", test.expectedOwner, owner)
			}
			if repo != test.expectedRepo {
				t.Errorf("Expected repo: %s, but got: %s", test.expectedRepo, repo)
			}
			if branch != test.expectedBranch {
				t.Errorf("Expected branch: %s, but got: %s", test.expectedBranch, branch)
			}
			if path != test.expectedPath {
				t.Errorf("Expected path: %s, but got: %s", test.expectedPath, path)
			}
		}
	}
}
