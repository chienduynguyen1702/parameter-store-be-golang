// routes/router.go

package routes

import (
	docs "parameter-store-be/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupV1Router() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	v1 := r.Group("/api/v1")
	{
		setupGroupAuth(v1)
		setupGroupOrganization(v1)
		setupGroupProject(v1)
		setupGroupSetting(v1)
	}
	// Swagger setup
	if gin.Mode() == gin.DebugMode {
		docs.SwaggerInfo.Title = "Parameter Store Backend API"
		docs.SwaggerInfo.Description = "This is a simple API for Parameter Store Backend."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = "localhost:8080"
		// docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else if gin.Mode() == gin.ReleaseMode {
		docs.SwaggerInfo.Title = "Parameter Store Backend API"
		docs.SwaggerInfo.Description = "This is a simple API for Parameter Store Backend."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = "parameter-store-be-golang.up.railway.app"
		// docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
