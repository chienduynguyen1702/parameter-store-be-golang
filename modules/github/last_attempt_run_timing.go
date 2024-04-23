package github

import (
	"encoding/json"
	"fmt"
	"io"
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
	Duration     int       `json:"duration"`
}

func makeGetWorkflowRunUsage(repoOwner string, repoName string, runID int, attemptNumber int, apiToken string) (*http.Request, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/attempts/%d", GitHubAPIEndpoint, repoOwner, repoName, runID, attemptNumber)
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

func GetLastAttemptNumberOfWorkflowRun(repoOwner string, repoName string, apiToken string, workflowName string) (int, int, error) {
	latestWorkflowID, statusCode, err := getWorkflowRunID(repoOwner, repoName, workflowName, apiToken)
	if err != nil {
		return http.StatusInternalServerError, 0, fmt.Errorf("error getting workflow run ID: %v", err)
	}
	if statusCode != http.StatusOK {
		return statusCode, 0, nil
	}

	client := &http.Client{}
	req, err := makeGetAWorkflowRuns(repoOwner, repoName, latestWorkflowID, apiToken)
	if err != nil {
		return http.StatusInternalServerError, 0, fmt.Errorf("error creating request: %v", err)
	}
	response, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, 0, fmt.Errorf("error sending request: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, 0, fmt.Errorf("error reading response: %v", err)
	}
	var run WorkflowRunResponse
	if err := json.Unmarshal(responseBody, &run); err != nil {
		return http.StatusInternalServerError, 0, fmt.Errorf("error unmarshalling response: %v", err)
	}
	// log.Println(run)
	return response.StatusCode, run.RunAttempt, nil
}

func GetWorkflowRunLastAttemptTimeDuration(repoOwner string, repoName string, runID int, attemptNumber int, apiToken string) (int, WorkflowRunAttempt, error) {
	client := &http.Client{}
	req, err := makeGetWorkflowRunUsage(repoOwner, repoName, runID, attemptNumber, apiToken)
	if err != nil {
		return http.StatusInternalServerError, WorkflowRunAttempt{}, fmt.Errorf("error creating request: %v", err)
	}
	response, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, WorkflowRunAttempt{}, fmt.Errorf("error sending request: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, WorkflowRunAttempt{}, fmt.Errorf("error reading response: %v", err)
	}
	var attempt WorkflowRunAttempt
	if err := json.Unmarshal(responseBody, &attempt); err != nil {
		return http.StatusInternalServerError, WorkflowRunAttempt{}, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return response.StatusCode, attempt, nil
}
