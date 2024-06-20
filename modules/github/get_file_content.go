package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FileContent represents the structure of the file content from GitHub
type FileContent struct {
	Type        string `json:"type"`
	Encoding    string `json:"encoding"`
	Size        int    `json:"size"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Content     string `json:"content"`
	SHA         string `json:"sha"`
	URL         string `json:"url"`
	GitURL      string `json:"git_url"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url"`
	Links       Links  `json:"_links"`
}

// Links represents the structure of the links within the file content
type Links struct {
	Git  string `json:"git"`
	Self string `json:"self"`
	HTML string `json:"html"`
}

func GetFileContent(owner, repo, path, token string) (string, error) {
	// create the request
	request, err := makeGetFileContentRequest(owner, repo, path, token)
	if err != nil {
		return "", err
	}
	// send the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("error sending request: ", err)
		return "", err
	}
	defer response.Body.Close()

	// read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// unmarshal the response body
	var fileContent FileContent
	if err := json.Unmarshal(responseBody, &fileContent); err != nil {
		return "", err
	}
	// return the content
	decodedContent, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return "", err
	}

	fmt.Println("decodedContent: \n", string(decodedContent))
	return string(decodedContent), nil
}

func makeGetFileContentRequest(owner, repo, path, token string) (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+repo+"/contents/"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
