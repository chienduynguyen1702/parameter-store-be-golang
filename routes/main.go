// routes/router.go

package routes

import (
	"os"
	docs "parameter-store-be/docs"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupV1Router() *gin.Engine {
	r := gin.Default()
	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://parameter-store-fe-golang.up.railway.app", "http://103.166.185.48:6888", "http://localhost:" + os.Getenv("PORT")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 30 * time.Hour,
	}))

	// Setup routes for the API version 1
	v1 := r.Group("/api/v1")
	{
		setupGroupAuth(v1)
		setupGroupOrganization(v1)
		setupGroupProjectList(v1)
		setupGroupProject(v1)
		setupGroupSetting(v1)
		setupGroupAgent(v1)
		setupGroupStage(v1)
		setupGroupEnvironment(v1)
	}
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	// Swagger setup
	docs.SwaggerInfo.Title = "Parameter Store Backend API"
	docs.SwaggerInfo.Description = "This is a simple API for Parameter Store Backend."
	docs.SwaggerInfo.Version = "1.0"
	if os.Getenv("ENVIRONMENT") == "dev" {
		docs.SwaggerInfo.Host = "localhost:" + os.Getenv("PORT")
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	} else if os.Getenv("ENVIRONMENT") == "production" {
		docs.SwaggerInfo.Host = "parameter-store-be-golang.up.railway.app"
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else if os.Getenv("ENVIRONMENT") == "datn-server" {
		docs.SwaggerInfo.Host = os.Getenv("HOSTNAME")
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
