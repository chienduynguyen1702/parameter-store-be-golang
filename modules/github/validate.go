package github

import (
	"fmt"
	"net/http"
)

func ValidateGithubRepo(repoURL string, repoAPIToken string) error {
	repo, err := ParseRepoURL(repoURL)
	if err != nil {
		return err
	}
	m, err := makeGetRepoInformationRequest(repo, repoAPIToken)
	if err != nil {
		return err
	}
	client := &http.Client{}
	response, err := client.Do(m)
	if err != nil {
		return fmt.Errorf("error sending request to get repo information: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthenticated by token to repo %s", repoURL)
	}
	if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("error to find repo %s", repoURL)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error getting repo information: %s", response.Status)
	}

	return nil
}

func makeGetRepoInformationRequest(repo Repository, repoAPIToken string) (*http.Request, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", GitHubAPIEndpoint, repo.Owner, repo.Name)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to get repo information: %v", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repoAPIToken))
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return request, nil
}

func ValidateWorkflowName(workflowName string, repoUrl string, apiToken string) error {
	repo, err := ParseRepoURL(repoUrl)
	if err != nil {
		return err
	}
	// debug
	// fmt.Println("repo owner: ", repo.Owner)
	// fmt.Println("repo name: ", repo.Name)
	// fmt.Println("apiToken: ", apiToken)
	idMatchedWorkflow, statusCode, err := getWorkflowID(repo.Owner, repo.Name, workflowName, apiToken)
	if err != nil {
		return err
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("Error getting workflow : %d", statusCode)
	}
	if idMatchedWorkflow == "" {
		return fmt.Errorf("Not found workflow name \"%s\" in github.com/%s/%s", workflowName, repo.Owner, repo.Name)
	}
	return nil
}
