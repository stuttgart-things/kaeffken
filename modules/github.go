/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"fmt"
	"strings"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/google/go-github/v62/github"
)

func CreateGithubClient(token string) *github.Client {
	if token == "" {
		log.Fatal("UNAUTHORIZED: NO TOKEN PRESENT")
	}
	return github.NewClient(nil).WithAuthToken(token)
}

func GetFileContentFromFileInGitHubRepo(client *github.Client, reference string) string {

	owner, repo, branch, path, _ := ParseGitHubURL(reference)

	fileContent, err := sthingsCli.GetFileContentFromGithubRepo(client, owner, repo, branch, path)
	if err != nil {
		log.Error("Error reading ", fileContent)
	}
	return fileContent

}

func ParseGitHubURL(url string) (owner, repo, branch, path string, err error) {
	// First, split the URL by "@" to separate the repo+branch part from the path
	parts := strings.Split(url, "@")
	if len(parts) != 2 {
		return "", "", "", "", fmt.Errorf("invalid format, missing '@'")
	}

	// Split the second part by ":" to separate the branch from the path
	branchPath := strings.Split(parts[1], ":")
	if len(branchPath) != 2 {
		return "", "", "", "", fmt.Errorf("invalid format, missing ':'")
	}
	branch = branchPath[0]
	path = branchPath[1]

	// Now, process the first part (before "@") to extract owner and repo
	repoPart := strings.TrimSuffix(parts[0], ".gi") // Strip ".gi" from the repo part
	ownerRepo := strings.Split(repoPart, "/")
	if len(ownerRepo) < 2 {
		return "", "", "", "", fmt.Errorf("invalid owner/repo format")
	}
	owner = ownerRepo[len(ownerRepo)-2]
	repo = ownerRepo[len(ownerRepo)-1]

	return owner, repo, branch, path, nil
}
