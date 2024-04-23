package github

import (
	"fmt"
	"net/http"
)

/*
curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR-TOKEN>" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/OWNER/REPO/actions/runs/RUN_ID/timing
*/

func makeGetWorkflowRunUsage(repoOwner string, repoName string, runID int, apiToken string) (*http.Request, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/timing", GitHubAPIEndpoint, repoOwner, repoName, runID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return req, nil
}
