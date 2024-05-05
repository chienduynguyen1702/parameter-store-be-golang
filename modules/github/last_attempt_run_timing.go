package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

/*
	curl -L \
	  -H "Accept: application/vnd.github+json" \
	  -H "Authorization: Bearer <YOUR-TOKEN>" \
	  -H "X-GitHub-Api-Version: 2022-11-28" \
	  https://api.github.com/repos/OWNER/REPO/actions/runs/RUN_ID/timing
*/
type WorkflowRunAttempt struct {
	RunAttempt   int       `json:"run_attempt"`
	RunStartedAt time.Time `json:"run_started_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Status       string    `json:"status"`
}

func makeGetWorkflowRunWithAttempt(repoOwner string, repoName string, workflowRunID int, attemptNumber int, apiToken string) (*http.Request, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/attempts/%d", GitHubAPIEndpoint, repoOwner, repoName, workflowRunID, attemptNumber)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return req, nil
}
func makeGetAWorkflowRuns(repoOwner string, repoName string, workflowrunID string, apiToken string) (*http.Request, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%s", GitHubAPIEndpoint, repoOwner, repoName, workflowrunID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return req, nil
}

type WorkflowRunResponse struct {
	ID         int       `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	RunAttempt int       `json:"run_attempt"`
}

func GetLastAttemptNumberOfWorkflowRun(repoOwner string, repoName string, apiToken string, workflowName string) (string, int, int, error) {
	latestWorkflowRunID, statusCode, err := getWorkflowRunID(repoOwner, repoName, workflowName, apiToken)
	if err != nil {
		return "", http.StatusInternalServerError, 0, fmt.Errorf("error getting workflow run ID: %v", err)
	}
	if statusCode != http.StatusOK {
		return "", statusCode, 0, nil
	}

	client := &http.Client{}
	req, err := makeGetAWorkflowRuns(repoOwner, repoName, latestWorkflowRunID, apiToken)
	if err != nil {
		return "", http.StatusInternalServerError, 0, fmt.Errorf("error creating request: %v", err)
	}
	response, err := client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, 0, fmt.Errorf("error sending request: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", http.StatusInternalServerError, 0, fmt.Errorf("error reading response: %v", err)
	}
	var run WorkflowRunResponse
	if err := json.Unmarshal(responseBody, &run); err != nil {
		return "", http.StatusInternalServerError, 0, fmt.Errorf("error unmarshalling response: %v", err)
	}
	// log.Println(run)
	return latestWorkflowRunID, response.StatusCode, run.RunAttempt, nil
}

func GetLastAttemptInformationOfWorkflowRun(repoOwner string, repoName string, apiToken string, workflowRunID int, attemptNumber int) (time.Time, time.Duration, error) {
	client := &http.Client{}
	req, err := makeGetWorkflowRunWithAttempt(repoOwner, repoName, workflowRunID, attemptNumber, apiToken)
	if err != nil {
		return time.Time{}, time.Duration(0), fmt.Errorf("error creating request: %v", err)
	}
	response, err := client.Do(req)
	if err != nil {
		return time.Time{}, time.Duration(0), fmt.Errorf("error sending request: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return time.Time{}, time.Duration(0), fmt.Errorf("error reading response: %v", err)
	}
	log.Println(string(responseBody))
	var run WorkflowRunAttempt
	if err := json.Unmarshal(responseBody, &run); err != nil {
		return time.Time{}, time.Duration(0), fmt.Errorf("error unmarshalling response: %v", err)
	}
	// log.Println(run)
	if run.Status != "completed" {
		return time.Time{}, time.Duration(0), fmt.Errorf("workflow run is not completed")
	}
	// return subtract of run.UpdateAt - run.RunStartedAt
	return run.RunStartedAt, run.UpdatedAt.Sub(run.RunStartedAt), nil

}
