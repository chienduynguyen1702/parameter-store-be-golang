package github

import (
	"fmt"
	"strings"
)

type Repository struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

// repoURL = "github.com/OWNER/REPO"
func ParseRepoURL(repoURL string) (Repository, error) {
	// split repoURL by "/"
	if repoURL == "" {
		return Repository{}, fmt.Errorf("repo URL is empty")
	}
	parts := strings.Split(repoURL, "/")
	if len(parts) != 3 {
		return Repository{}, fmt.Errorf("invalid repo URL: %s", repoURL)
	}
	var g Repository
	g.Owner = parts[1]
	g.Name = parts[2]
	return g, nil
}
