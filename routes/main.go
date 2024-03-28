// routes/router.go

package routes

import (
	docs "parameter-store-be/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupV1Router() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		setupGroupAuth(v1)
		setupGroupOrganization(v1)
		setupGroupProjectOverview(v1)
		setupGroupProjectParameter(v1)
		setupGroupProjectAgent(v1)
		setupGroupProjectTracking(v1)
		// setupGroupUserSetting(v1)
	}
	// Swagger setup
	docs.SwaggerInfo.Title = "Parameter Store Backend API"
	docs.SwaggerInfo.Description = "This is a simple API for Parameter Store Backend."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "parameter-store-be-golang.up.railway.com"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
