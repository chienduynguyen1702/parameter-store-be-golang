package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserRepositorys struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
	Private bool   `json:"private"`
	HtmlUrl string `json:"html_url"`
}

func GetUserRepos(username, token string) ([]UserRepositorys, error) {
	var repos []UserRepositorys
	perPage := 100
	for page := 1; ; page++ {
		tempRepos, err := getUserReposRequest(username, token, page, perPage)
		if err != nil {
			return nil, err
		}
		if tempRepos == nil {
			break
		}
		repos = append(repos, tempRepos...)
		if len(tempRepos) < perPage {
			break
		}

	}
	return repos, nil
}

func getUserReposRequest(username, token string, page, perPage int) ([]UserRepositorys, error) {
	req, err := makeGetUserReposRequest(username, token, page, perPage)
	if err != nil {
		return nil, err
	}

	// send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse the response
	var repos []UserRepositorys
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func makeGetUserReposRequest(username, token string, page, perPage int) (*http.Request, error) {
	//		curl -L \
	//	  -H "Accept: application/vnd.github+json" \
	//	  -H "Authorization: Bearer <token>" \
	//	  -H "X-GitHub-Api-Version: 2022-11-28" \
	//	  https://api.github.com/user/repos?sort=pushed
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/user/repos?sort=pushed&page=%d&per_page=%d&type=all", page, perPage), nil)
	if err != nil {
		return nil, err
	}

	// set the request headers
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return request, nil
}
