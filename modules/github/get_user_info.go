package github

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Plan represents the plan details in the JSON response
type Plan struct {
	Name          string `json:"name"`
	Space         int    `json:"space"`
	PrivateRepos  int    `json:"private_repos"`
	Collaborators int    `json:"collaborators"`
}

// GitHubUser represents the user details in the JSON response
type GitHubUser struct {
	Login             string    `json:"login"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PrivateGists      int       `json:"private_gists"`
	TotalPrivateRepos int       `json:"total_private_repos"`
	OwnedPrivateRepos int       `json:"owned_private_repos"`
	DiskUsage         int       `json:"disk_usage"`
	Collaborators     int       `json:"collaborators"`
	TwoFactorAuth     bool      `json:"two_factor_authentication"`
	Plan              Plan      `json:"plan"`
}

func makeGetGitUserInfo(token string) (*http.Request, error) {
	// create the request
	request, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	// set the request headers
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "token "+token)
	return request, nil
}

func getGitUserInfoRequest(token string) (*http.Response, error) {
	// create the request
	request, err := makeGetGitUserInfo(token)
	if err != nil {
		return nil, err
	}

	// send the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func GetGitUserInfo(token string) (GitHubUser, error) {
	// create the request
	response, err := getGitUserInfoRequest(token)
	if err != nil {
		return GitHubUser{}, err
	}
	var userInfo GitHubUser
	// read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return GitHubUser{}, err
	}

	// unmarshal the response body
	if err := json.Unmarshal(responseBody, &userInfo); err != nil {
		return GitHubUser{}, err
	}

	return userInfo, nil
}
