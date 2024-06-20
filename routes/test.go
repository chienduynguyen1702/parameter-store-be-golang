package routes

import (
	"parameter-store-be/controllers"

	"github.com/gin-gonic/gin"
)

func setupGroupTestAPI(r *gin.RouterGroup) {
	testGithubGroup := r.Group("/test")
	{
		testGithubGroup.PUT("/update-secrets", controllers.TestUpdateSecrets)
		testGithubGroup.PUT("/get-file-content", controllers.TestGetFileContent)
	}
}
