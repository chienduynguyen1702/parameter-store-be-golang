package controllers

import (
	"parameter-store-be/modules/github"

	"github.com/gin-gonic/gin"
)

type updateSecretBody struct {
	SecretName     string `json:"secret_name"`
	Value          string `json:"value"`
	Owner          string `json:"owner"`
	Repo           string `json:"repo"`
	Token          string `json:"token"`
	EncryptedValue string `json:"encrypted_value"`
}

// UpdateSecrets godoc
// @Summary Update a secret in a github repository
// @Description Update a secret in a github repository
// @ID update-secrets
// @Tags Test / Github
// @Accept  json
// @Produce  json
// @Param body body controllers.updateSecretBody true "Update a secret in a github repository"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/test/update-secrets [put]
func TestUpdateSecrets(c *gin.Context) {
	var body updateSecretBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := github.CreateSecrets(body.Owner, body.Repo, body.SecretName, body.Value, body.Token)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// return the response body
	c.JSON(200, gin.H{"message": "Secret updated successfully"})

}

type getFileContentBody struct {
	Path  string `json:"path"`
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Token string `json:"token"`
}

// GetFileContent godoc
// @Summary Get the content of a file in a github repository
// @Description Get the content of a file in a github repository
// @ID get-file-content
// @Tags Test / Github
// @Accept  json
// @Produce  json
// @Param body body controllers.getFileContentBody true "Get the content of a file in a github repository"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/test/get-file-content [put]
func TestGetFileContent(c *gin.Context) {
	var body getFileContentBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	content, err := github.GetFileContent(body.Owner, body.Repo, body.Path, body.Token)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// return the response body
	c.JSON(200, gin.H{"content": content})
}
