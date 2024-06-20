package github

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/crypto/nacl/box"
)

// const keyID = "0123456789!@#$%^&*()"

func CreateSecrets(owner, repo, secretName, value, token string) error {
	// get public key id
	keyID, keyValue, err := GetPublicKeyIDOfRepo(owner, repo, token)
	if err != nil {
		return err
	}

	//encrypt the value
	encryptedValue, err := EncryptedValue(value, keyValue)
	if err != nil {
		return err
	}

	// make the http client
	req, err := MakeCreateSecretsHTTPClient(owner, repo, secretName, encryptedValue, token, keyID)
	if err != nil {
		return err
	}
	// make the request by passing the client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// // read the response body
	// responseBody, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("response status code", resp.StatusCode)
	// fmt.Println("responseBody", string(responseBody))
	status := resp.StatusCode
	if status == 201 {
		fmt.Println("Secret created successfully")
	} else if status == 204 {
		fmt.Println("Secret updated successfully")
	} else {
		fmt.Println("Error while creating secret")
	}

	return nil
}

func MakeCreateSecretsHTTPClient(owner, repo, secretName string, encryptedValue string, token string, keyID string) (*http.Request, error) {
	//		curl -L \
	//	  -X PUT \
	//	  -H "Accept: application/vnd.github+json" \
	//	  -H "Authorization: Bearer <YOUR-TOKEN>" \
	//	  -H "X-GitHub-Api-Version: 2022-11-28" \
	//	  https://api.github.com/repos/OWNER/REPO/actions/secrets/SECRET_NAME \
	//	  -d '{"encrypted_value":"c2VjcmV0","key_id":"012345678912345678"}'
	req, err := http.NewRequest("PUT", "https://api.github.com/repos/"+owner+"/"+repo+"/actions/secrets/"+secretName, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	req.Header.Set("Content-Type", "application/json")

	// set the body
	body := `{"encrypted_value":"` + encryptedValue + `","key_id":"` + keyID + `"}`
	req.Body = io.NopCloser(strings.NewReader(body))

	return req, nil
}

func EncryptedValue(data, publicKey string) (string, error) {
	// https://stackoverflow.com/questions/76562205/how-to-encrypt-repository-secret-for-github-action-secrets-api
	// Decode the public key from base64
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	// Decode the public key
	var publicKeyDecoded [32]byte
	copy(publicKeyDecoded[:], publicKeyBytes)

	// Encrypt the secret value
	encrypted, err := box.SealAnonymous(nil, []byte(data), (*[32]byte)(publicKeyBytes), rand.Reader)

	if err != nil {
		return "", err
	}
	// Encode the encrypted value in base64
	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)
	// fmt.Println("encryptedBase64 :", encryptedBase64)
	return encryptedBase64, nil

}

type PublicKey struct {
	KeyID string `json:"key_id"`
	Key   string `json:"key"`
}

func GetPublicKeyIDOfRepo(onwer, repo, token string) (string, string, error) {
	// curl -L \
	//   -H "Accept: application/vnd.github.v3+json" \
	//   -H "Authorization`

	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+onwer+"/"+repo+"/actions/secrets/public-key", nil)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	defer resp.Body.Close()
	// read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response Body", err)
		return "", "", err
	}

	// fmt.Println(string(responseBody))

	// bind to struct
	var publicKey PublicKey
	err = json.Unmarshal(responseBody, &publicKey)
	if err != nil {
		fmt.Println("Error while unmarshalling response", err)
		return "", "", err
	}

	// fmt.Println("publicKey.KeyID :", publicKey.KeyID)
	// fmt.Println("publicKey.Key   :", publicKey.Key)
	return publicKey.KeyID, publicKey.Key, nil
}
