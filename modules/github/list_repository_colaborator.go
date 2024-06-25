package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RepositoryColaborator struct {
	Login       string `json:"login"`
	Permissions struct {
		Admin bool `json:"admin"`
	} `json:"permissions"`
	Email string `json:"email"`
}

func ListRepositoryColaborator(owner, repository, token string) ([]RepositoryColaborator, error) {
	var colaborators []RepositoryColaborator
	perPage := 100
	for page := 1; ; page++ {
		tempColaborators, err := listRepositoryColaboratorRequest(owner, repository, token, page, perPage)
		if err != nil {
			return nil, err
		}
		if tempColaborators == nil {
			break
		}
		colaborators = append(colaborators, tempColaborators...)
		if len(tempColaborators) < perPage {
			break
		}

	}
	return colaborators, nil
}

func listRepositoryColaboratorRequest(owner, repository, token string, page, perPage int) ([]RepositoryColaborator, error) {
	req, err := makeListRepositoryColaboratorRequest(owner, repository, token, page, perPage)
	if err != nil {
		return nil, err
	}

	// send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request to Github API")
		return nil, err
	}
	defer resp.Body.Close()

	// parse the response
	if resp.StatusCode != http.StatusOK {
		log.Println("Error response from Github API")
		return nil, nil
	}
	// parse response to string
	// stringBody, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println("Error reading response from Github API")
	// 	return nil, err
	// }
	// log.Println(string(stringBody))
	var colaborators []RepositoryColaborator
	err = json.NewDecoder(resp.Body).Decode(&colaborators)
	if err != nil {
		log.Println("Error parsing response from Github API")
		return nil, err
	}
	return colaborators, nil
}

func makeListRepositoryColaboratorRequest(owner, repository, token string, page, perPage int) (*http.Request, error) {
	// create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/collaborators?page=%d&per_page=%d", owner, repository, page, perPage), nil)
	if err != nil {
		return nil, err
	}

	// set the request headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}
