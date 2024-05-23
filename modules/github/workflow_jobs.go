package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WorkflowRunJobs struct {
	Total int    `json:"total_count"`
	Jobs  []Jobs `json:"jobs"`
}

func (w WorkflowRunJobs) Print() {
	fmt.Printf("Total: %d, Jobs: %v", w.Total, w.Jobs)
	for _, job := range w.Jobs {
		fmt.Printf("===================JOB==================\n")
		fmt.Printf("ID: %d, RunID: %d, Name: %s, Status: %s, Conclusion: %s\n", job.ID, job.RunID, job.Name, job.Status, job.Conclusion)
		for _, step := range job.Steps {
			fmt.Printf("\tNumber: %d, Conclusion: %s, Name: %s, Status: %s\n", step.Number, step.Conclusion, step.Name, step.Status)
		}
	}
}

type Jobs struct {
	ID         int    `json:"id"`
	RunID      int    `json:"run_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	Steps      []Step `json:"steps"`
}

type Step struct {
	Number     int    `json:"number"`
	Conclusion string `json:"conclusion"`
	Name       string `json:"name"`
	Status     string `json:"status"`
}

func ListJobsForAWorkflowRun(owner string, repo string, token string, runID int) (WorkflowRunJobs, error) {
	var workflowRunJobs WorkflowRunJobs
	client := &http.Client{}
	request, err := makeListJobsForAWorkflowRun(owner, repo, token, runID)
	if err != nil {
		return workflowRunJobs, fmt.Errorf("error creating request to get jobs for a workflow run: %v", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return workflowRunJobs, fmt.Errorf("error making request to get jobs for a workflow run: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		// fmt.Println("Error reading response get workflow ID:", err)
		return WorkflowRunJobs{}, fmt.Errorf("error reading response get workflow ID: %v", err)
	}
	// unmarschal the response body
	err = json.Unmarshal(responseBody, &workflowRunJobs)

	return workflowRunJobs, err
}

func makeListJobsForAWorkflowRun(owner string, repo string, repoAPIToken string, runID int) (*http.Request, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/jobs", GitHubAPIEndpoint, owner, repo, runID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to get repo information: %v", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repoAPIToken))
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return request, nil
}
