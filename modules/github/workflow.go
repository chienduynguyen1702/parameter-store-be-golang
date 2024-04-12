package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	GitHubAPIEndpoint = "https://api.github.com"
)

type WorkflowRunsResponse struct {
	TotalCount   int `json:"total_count"`
	WorkflowRuns []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"workflow_runs"`
}

/*
input
- repoOwner : owner of the repository
- repoName : name of the repository
- workflowName : name of the workflow
- apiToken : github api token
output
- status : if Rerun is successful or not
- error : error if any
*/
func RerunWorkFlow(repoOwner string, repoName string, workflowName string, apiToken string) (int, error) {
	// API docs : https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#re-run-a-workflow
	// curl -L \
	// -X POST \
	// -H "Accept: application/vnd.github+json" \
	// -H "Authorization: Bearer <YOUR-TOKEN>" \
	// -H "X-GitHub-Api-Version: 2022-11-28" \
	// https://api.github.com/repos/OWNER/REPO/actions/runs/RUN_ID/rerun

	// 1 - Get the latest workflow run
	workflowID, statusCode, err := getWorkflowID(repoOwner, repoName, workflowName, apiToken)
	if err != nil {
		return statusCode, err
	}
	// 2 - Rerun the workflow
	client := &http.Client{}
	rerunWorkflowRequest, err := makeNewRerunWorkflowRequest(repoOwner, repoName, workflowID, apiToken)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error creating request for rerun: %v", err)
	}
	// Send the request
	response, err := client.Do(rerunWorkflowRequest)
	if err != nil {
		return response.StatusCode, fmt.Errorf("error sending request for rerun: %v", err)
	}

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error reading response for rerun: %v", err)
	}

	// debug
	// fmt.Printf("Rerun workflow response: %s\n", string(responseBody))

	// if resonse status !=201 then return error
	if response.StatusCode != 201 {
		return response.StatusCode, fmt.Errorf("error rerunning workflow: %s ", string(responseBody))
	}

	return response.StatusCode, nil
}
func makeNewRerunWorkflowRequest(repoOwner string, repoName string, workflowID string, apiToken string) (*http.Request, error) {
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

func getWorkflowID(repoOwner string, repoName string, workflowName string, apiToken string) (string, int, error) {
	// API docs : https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#list-workflow-runs-for-a-repository
	// curl -L \
	// -H "Accept: application/vnd.github+json" \
	// -H "Authorization: Bearer <YOUR-TOKEN>" \
	// -H "X-GitHub-Api-Version: 2022-11-28" \
	// https://api.github.com/repos/OWNER/REPO/actions/runs

	// Create a new HTTP client
	client := &http.Client{}
	listWorkflowRequest, err := makeNewListWorkflowRequest(repoOwner, repoName, apiToken)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error creating request to get workflow ID: %v", err)
	}
	// Send the GET request
	response, err := client.Do(listWorkflowRequest)
	if err != nil {
		// fmt.Println("Error sending request to get workflow ID:", err)
		return "", response.StatusCode, fmt.Errorf("error sending request to get workflow ID: %v", err)
	}
	defer response.Body.Close()

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

	// debug
	// fmt.Println("List workflow: ", workflowRunsResponse)

	// Find the workflow ID matching the workflow name
	var idMatchedWorkflow string
	workflowIsFound := false
	for _, workflowRun := range workflowRunsResponse.WorkflowRuns {
		if workflowRun.Name == workflowName {
			// fmt.Printf("Matched workflow ID: \"%d\"\n", workflowRun.ID)
			// fmt.Printf("Matched workflow Name: \"%s\"\n", workflowRun.Name)
			idMatchedWorkflow = fmt.Sprintf("%d", workflowRun.ID)
			workflowIsFound = true
			break
		}
	}
	if !workflowIsFound {
		return "", http.StatusNotFound, fmt.Errorf("workflow name \"%s\" not found", workflowName)
	}

	return idMatchedWorkflow, response.StatusCode, nil
}

func makeNewListWorkflowRequest(repoOwner string, repoName string, apiToken string) (*http.Request, error) {

	// Create a new GET request
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs", GitHubAPIEndpoint, repoOwner, repoName)
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
