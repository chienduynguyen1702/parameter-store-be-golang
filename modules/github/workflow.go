package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	GitHubAPIEndpoint = "https://api.github.com"
	MAX_PAGE_SIZE     = 100
	DEFAULT_PAGE_SIZE = 30
)

type WorkflowRunsResponse struct {
	TotalCount   int            `json:"total_count"`
	WorkflowRuns []WorkflowRuns `json:"workflow_runs"`
}
type WorkflowRuns struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	RunAttempt   int    `json:"run_attempt"`
	DisplayTitle string `json:"display_title"`
	CreatedAt    string `json:"created_at"`
}

/*
input
- repoOwner : owner of the repository
- repoName : name of the repository
- workflowName : name of the workflow
- apiToken : github api token
output
- response message: message of the action rerun
- response status : status of the action rerun
- error : error if any
*/
func RerunWorkFlow(repoOwner string, repoName string, workflowName string, apiToken string) (int, string, error) {
	// API docs : https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#re-run-a-workflow
	// curl -L \
	// -X POST \
	// -H "Accept: application/vnd.github+json" \
	// -H "Authorization: Bearer <YOUR-TOKEN>" \
	// -H "X-GitHub-Api-Version: 2022-11-28" \
	// https://api.github.com/repos/OWNER/REPO/actions/runs/RUN_ID/rerun

	// 1 - Get the latest workflow run
	workflowID, statusCode, err := getWorkflowRunID(repoOwner, repoName, workflowName, apiToken)
	if err != nil {
		return statusCode, err.Error(), err
	}
	// 2 - Rerun the workflow
	client := &http.Client{}
	rerunWorkflowRequest, err := makeRerunWorkflowRequest(repoOwner, repoName, workflowID, apiToken)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("error creating request for rerun: %v", err)
	}
	// Send the request
	response, errReq := client.Do(rerunWorkflowRequest)

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("error reading response for rerun: %v", err)
	}

	return response.StatusCode, string(responseBody), errReq
}
func makeRerunWorkflowRequest(repoOwner string, repoName string, workflowID string, apiToken string) (*http.Request, error) {
	// Create a new POST request
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%s/rerun", GitHubAPIEndpoint, repoOwner, repoName, workflowID)
	// fmt.Println("NewRerunWorkflowRequest URL: ", url)

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		// fmt.Println("Error creating request:", err)
		return nil, fmt.Errorf("error creating rerun workflow request: %v", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	// debug
	// fmt.Println("NewRerunWorkflowRequest:", request)
	return request, nil
}

func getWorkflowRunID(repoOwner string, repoName string, workflowName string, apiToken string) (string, int, error) {
	// API docs : https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#list-workflow-runs-for-a-repository
	// curl -L \
	// -H "Accept: application/vnd.github+json" \
	// -H "Authorization: Bearer <YOUR-TOKEN>" \
	// -H "X-GitHub-Api-Version: 2022-11-28" \
	// https://api.github.com/repos/OWNER/REPO/actions/runs

	// Create a new HTTP client
	client := &http.Client{}
	listWorkflowRequest, err := makeListWorkflowRunRequest(repoOwner, repoName, apiToken)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error creating request to get workflow ID: %v", err)
	}
	// Send the GET request
	response, err := client.Do(listWorkflowRequest)
	if err != nil {
		fmt.Println("Error sending request to get workflow ID:", err)
		return "", response.StatusCode, fmt.Errorf("error sending request to get workflow ID: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized {
		// fmt.Println("Error response get workflow ID:", response.Status)
		return "", response.StatusCode, fmt.Errorf("unauthenticated by token to repo github.com/%v/%v", repoOwner, repoName)
	}
	if response.StatusCode == http.StatusNotFound {
		// fmt.Println("Error response get workflow ID:", response.Status)
		return "", response.StatusCode, fmt.Errorf("error to find repo github.com/%v/%v", repoOwner, repoName)
	}
	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		// fmt.Println("Error reading response get workflow ID:", err)
		return "", http.StatusInternalServerError, fmt.Errorf("error reading response get workflow ID: %v", err)
	}
	// Unmarshal the response body
	var workflowRunsResponse WorkflowRunsResponse
	err = json.Unmarshal(responseBody, &workflowRunsResponse)
	if err != nil {
		// fmt.Println("Error unmarshalling response get workflow ID:", err)
		return "", http.StatusInternalServerError, fmt.Errorf("error unmarshalling response get workflow ID: %v", err)
	}

	// Find the workflow ID matching the workflow name
	var idMatchedWorkflow string
	workflowIsFound := false

	// get the total workflow runs, then check if total workflow runs is bigger than MAX_PAGE_SIZE
	totalWorkflowRuns := workflowRunsResponse.TotalCount
	// check if first workflow run list is contain the workflow name
	for _, workflowRun := range workflowRunsResponse.WorkflowRuns {

		if workflowRun.Name == workflowName {
			// fmt.Printf("Matched workflow ID: \"%d\"\n", workflowRun.ID)
			// fmt.Printf("Matched workflow Name: \"%s\"\n", workflowRun.Name)
			idMatchedWorkflow = fmt.Sprintf("%d", workflowRun.ID)
			workflowIsFound = true
			break
		}
	}
	listWorkflowRuns := make([]WorkflowRuns, 0)
	if totalWorkflowRuns > MAX_PAGE_SIZE && !workflowIsFound {
		// Allocate memory for WorkflowRuns slice of struct with totalWorkflowRuns

		// Calculate total pages
		totalPages := (totalWorkflowRuns-1)/MAX_PAGE_SIZE + 1

		// Create a new HTTP client
		client := &http.Client{}

		for page := 2; page <= totalPages; page++ {
			listWorkflowRequest, err := makeListWorkflowRunRequestWithPage(repoOwner, repoName, apiToken, page)
			if err != nil {
				return "", http.StatusInternalServerError, fmt.Errorf("error creating request to get workflow ID: %v", err)
			}

			// Send the GET request
			response, err := client.Do(listWorkflowRequest)
			if err != nil {
				fmt.Println("Error sending request to get workflow ID:", err)
				return "", response.StatusCode, fmt.Errorf("error sending request to get workflow ID: %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusUnauthorized {
				// fmt.Println("Error response get workflow ID:", response.Status)
				return "", response.StatusCode, fmt.Errorf("unauthenticated by token to repo github.com/%v/%v", repoOwner, repoName)
			}
			if response.StatusCode == http.StatusNotFound {
				// fmt.Println("Error response get workflow ID:", response.Status)
				return "", response.StatusCode, fmt.Errorf("error to find repo github.com/%v/%v", repoOwner, repoName)
			}

			// Read the response body
			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				// fmt.Println("Error reading response get workflow ID:", err)
				return "", http.StatusInternalServerError, fmt.Errorf("error reading response get workflow ID: %v", err)
			}

			var pageWorkflowRuns WorkflowRunsResponse
			err = json.Unmarshal(responseBody, &pageWorkflowRuns)
			if err != nil {
				return "", http.StatusInternalServerError, fmt.Errorf("error unmarshalling response get workflow ID: %v", err)
			}

			// Append the workflow runs to the list
			listWorkflowRuns = append(listWorkflowRuns, pageWorkflowRuns.WorkflowRuns...)
			for _, workflowRun := range listWorkflowRuns {
				// fmt.Printf("################ %d ##############", i)
				// fmt.Printf("Workflow ID: \"%d\"\n", workflowRun.ID)
				// fmt.Printf("Workflow Name: \"%s\"\n", workflowRun.Name)
				// fmt.Printf("Workflow CreatedAt: \"%s\"\n", workflowRun.CreatedAt)
				// fmt.Printf("Workflow DisplayTitle: \"%s\"\n", workflowRun.DisplayTitle)

				if workflowRun.Name == workflowName {
					// fmt.Printf("Matched workflow ID: \"%d\"\n", workflowRun.ID)
					// fmt.Printf("Matched workflow Name: \"%s\"\n", workflowRun.Name)
					idMatchedWorkflow = fmt.Sprintf("%d", workflowRun.ID)
					workflowIsFound = true
					break
				}
			}
		}

	}

	if !workflowIsFound {
		return "", http.StatusNotFound, fmt.Errorf("not found workflow name \"%s\" in github.com/%v/%v", workflowName, repoOwner, repoName)
	}

	return idMatchedWorkflow, response.StatusCode, nil
}

func makeListWorkflowRunRequest(repoOwner string, repoName string, apiToken string) (*http.Request, error) {

	// Create a new GET request
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs?per_page=%d", GitHubAPIEndpoint, repoOwner, repoName, MAX_PAGE_SIZE)
	// fmt.Println("ListWorkflowRequest URL: ", url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// fmt.Println("Error creating request:", err)
		return nil, fmt.Errorf("error creating list workflow request: %v", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	// debug
	// fmt.Println("ListWorkflowRequest: ", request)
	return request, nil
}
func makeListWorkflowRunRequestWithPage(repoOwner string, repoName string, apiToken string, page int) (*http.Request, error) {

	// Create a new GET request
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs?per_page=%d&page=%d", GitHubAPIEndpoint, repoOwner, repoName, MAX_PAGE_SIZE, page)
	// fmt.Println("ListWorkflowRequest URL: ", url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// fmt.Println("Error creating request:", err)
		return nil, fmt.Errorf("error creating list workflow request: %v", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	// debug
	// fmt.Println("ListWorkflowRequest: ", request)
	return request, nil
}

func makeGetWorkflowRunRequest(repoOwner string, repoName string, apiToken string, workflowID string) (*http.Request, error) {

	// Create a new GET request
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%s", GitHubAPIEndpoint, repoOwner, repoName, workflowID)
	// fmt.Println("ListWorkflowRequest URL: ", url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// fmt.Println("Error creating request:", err)
		return nil, fmt.Errorf("error creating list workflow request: %v", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	// debug
	// fmt.Println("ListWorkflowRequest: ", request)
	return request, nil
}

type WorkflowsResponse struct {
	TotalCount int `json:"total_count"`
	Workflows  []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Path  string `json:"path"`
		State string `json:"state"`
	} `json:"workflows"`
}

func (w WorkflowsResponse) Print() {
	fmt.Printf("TotalCount: %d\n\n", w.TotalCount)
	for _, workflow := range w.Workflows {
		fmt.Printf("ID: %d\n", workflow.ID)
		fmt.Printf("Name: %s\n", workflow.Name)
		fmt.Printf("Path: %s\n", workflow.Path)
		fmt.Printf("State: %s\n", workflow.State)
	}
}

func makeListWorkflowsRequest(repoOwner string, repoName string, apiToken string) (*http.Request, error) {

	// Create a new GET request
	url := fmt.Sprintf("%s/repos/%s/%s/actions/workflows", GitHubAPIEndpoint, repoOwner, repoName)
	// fmt.Println("ListWorkflowRequest URL: ", url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// fmt.Println("Error creating request:", err)
		return nil, fmt.Errorf("error creating list workflow request: %v", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	// debug
	// fmt.Println("ListWorkflowRequest: ", request)
	return request, nil
}
func GetWorkflows(RepoURL string, apiToken string) (WorkflowsResponse, error) {
	// Parse the repository URL
	repo, err := ParseRepoURL(RepoURL)
	if err != nil {
		return WorkflowsResponse{}, fmt.Errorf("error parsing repository URL: %v", err)
	}

	// Create a new HTTP client
	client := &http.Client{}
	listWorkflowRequest, err := makeListWorkflowsRequest(repo.Owner, repo.Name, apiToken)
	if err != nil {
		return WorkflowsResponse{}, fmt.Errorf("error creating request to get workflow ID: %v", err)
	}
	// Send the GET request
	response, err := client.Do(listWorkflowRequest)
	if err != nil {
		fmt.Println("Error sending request to get workflow ID:", err)
		return WorkflowsResponse{}, fmt.Errorf("error sending request to get workflow ID: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized {
		// fmt.Println("Error response get workflow ID:", response.Status)
		return WorkflowsResponse{}, fmt.Errorf("unauthenticated by token to repo github.com/%v/%v", repo.Owner, repo.Name)
	}
	if response.StatusCode == http.StatusNotFound {
		// fmt.Println("Error response get workflow ID:", response.Status)
		return WorkflowsResponse{}, fmt.Errorf("error to find repo github.com/%v/%v", repo.Owner, repo.Name)
	}
	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		// fmt.Println("Error reading response get workflow ID:", err)
		return WorkflowsResponse{}, fmt.Errorf("error reading response get workflow ID: %v", err)
	}
	// debug
	// fmt.Println("Response body: ", string(responseBody))
	// Unmarshal the response body
	var workflowsResponse WorkflowsResponse
	err = json.Unmarshal(responseBody, &workflowsResponse)
	if err != nil {
		// fmt.Println("Error unmarshalling response get workflow ID:", err)
		return WorkflowsResponse{}, fmt.Errorf("error unmarshalling response get workflow ID: %v", err)
	}
	// workflowsResponse.Print()
	// Remove the dynamic workflow file
	workflowsResponse.removeDynamicWorkflowfile()
	// debug
	// fmt.Println("List workflow: ", workflowsResponse)

	return workflowsResponse, nil
}

func (w *WorkflowsResponse) removeDynamicWorkflowfile() {
	// if path is dynamic/pages/pages-build-deployment
	// in the path, the first pages is dynamic
	// then remove this workflow
	for i, workflow := range w.Workflows {
		if isDynamicWorkflow(workflow.Path) {
			w.Workflows = append(w.Workflows[:i], w.Workflows[i+1:]...)
			break
		}
	}
}

func isDynamicWorkflow(workflowPath string) bool {
	// tokenize the path  by "/"
	// if first token is "dynamic"
	// then it is dynamic workflow dynamic/pages/pages-build-deployment
	workflowPathParts := splitWorkflowPath(workflowPath)
	if len(workflowPathParts) > 0 && workflowPathParts[0] == "dynamic" {
		return true
	}
	return false
}

func splitWorkflowPath(workflowPath string) []string {
	workflowPathParts := strings.Split(workflowPath, "/")
	// split the path by "/"
	return workflowPathParts
}
