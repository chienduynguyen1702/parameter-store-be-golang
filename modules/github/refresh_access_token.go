package github

// func GetAccessToken(code string) (string, error) {
// 	// create the request
// 	response, err := getAccessTokenRequest(code)
// 	if err != nil {
// 		return "", err
// 	}

// 	// read the response body
// 	responseBody, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	// parse the response
// 	var accessToken AccessToken
// 	err = json.Unmarshal(responseBody, &accessToken)
// 	if err != nil {
// 		return "", err
// 	}

// 	return accessToken.AccessToken, nil
// }

// func getAccessTokenRequest(code string) (*http.Response, error) {
// 	// create the request
// 	request, err := makeGetAccessTokenRequest(code)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// send the request
// 	response, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

// func makeGetAccessTokenRequest(code string) (*http.Request, error) {
// 	// create the request
// 	request, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&grant_type=refresh_token", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// set the request headers
// 	request.Header.Set("Accept", "application/json")
// 	request.Header.Set("Content-Type", "application/json")

// 	// set the request body
// 	requestBody := strings.NewReader(`{"client_id":"` + os.Getenv("GITHUB_CLIENT_ID") + `","client_secret":"` + os.Getenv("GITHUB_CLIENT_SECRET") + `","code":"` + code + `"}`)
// 	return request, nil
// }
