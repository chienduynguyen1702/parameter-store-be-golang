package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SearchCodeInRepoResponse struct {
	TotalCount        int                    `json:"total_count"`
	IncompleteResults bool                   `json:"incomplete_results"`
	Items             []SearchCodeInRepoItem `json:"items"`
}
type SearchCodeInRepoItem struct {
	Name    string  `json:"name"`
	Path    string  `json:"path"`
	SHA     string  `json:"sha"`
	HTMLURL string  `json:"html_url"`
	Score   float64 `json:"score"`
}

func makeSearchCodeInRepoRequest(owner, repo, token, searchStr string) (*http.Request, error) {
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s+in:file+repo:%s/%s", searchStr, owner, repo)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return request, nil
}

func searchCodeInRepoRequest(owner, repo, token, searchStr string) (*http.Response, error) {
	// create the request
	request, err := makeSearchCodeInRepoRequest(owner, repo, token, searchStr)
	if err != nil {
		return nil, err
	}
	// send the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("error sending request: ", err)
		return nil, err
	}
	return response, nil
}

func SearchCodeInRepo(owner, repo, token, searchStr string) (SearchCodeInRepoResponse, error) {
	// send the request
	response, err := searchCodeInRepoRequest(owner, repo, token, searchStr)
	if err != nil {
		log.Printf("error sending request: %v", err)
		return SearchCodeInRepoResponse{}, err
	}
	defer response.Body.Close()

	// read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
		return SearchCodeInRepoResponse{}, err
	}

	// unmarshal the response body
	var searchCodeInRepoResponse SearchCodeInRepoResponse
	if err := json.Unmarshal(responseBody, &searchCodeInRepoResponse); err != nil {
		log.Printf("error unmarshalling response body: %v", err)
		return SearchCodeInRepoResponse{}, err
	}

	return searchCodeInRepoResponse, nil
}
