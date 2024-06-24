package github

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func makeGetAccessTokenRequest(code string) (*http.Request, error) {
	// create the request
	request, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token"+"?client_id="+os.Getenv("GITHUB_CLIENT_ID")+"&client_secret="+os.Getenv("GITHUB_CLIENT_SECRET")+"&code="+code, nil)
	if err != nil {
		return nil, err
	}

	// set the request headers
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	// set the request body
	requestBody := strings.NewReader(`{"client_id":"` + os.Getenv("GITHUB_CLIENT_ID") + `","client_secret":"` + os.Getenv("GITHUB_CLIENT_SECRET") + `","code":"` + code + `"}`)
	request.Body = io.NopCloser(requestBody)
	return request, nil
}
func getAccessTokenRequest(code string) (*http.Response, error) {
	// create the request
	request, err := makeGetAccessTokenRequest(code)
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
func GetAccessToken(code string) (string, error) {
	// create the request
	response, err := getAccessTokenRequest(code)
	if err != nil {
		return "", err
	}

	// read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// unmarshal the response body
	var accessToken AccessToken
	if err := json.Unmarshal(responseBody, &accessToken); err != nil {
		return "", err
	}

	// return the access token
	return accessToken.AccessToken, nil
}
