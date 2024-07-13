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
	whiteList := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://192.168.88.153:3000",
		"https://parameter-store-fe-golang.up.railway.app",
		os.Getenv("HOSTNAME_URL"),
		os.Getenv("CLIENT_URL"),
		"https://param-store.datn.live",
		"https://chienduynguyen1702.github.io",
		"http://localhost:" + os.Getenv("PORT"),
	}
	r := gin.Default()
	
	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     whiteList,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Content-Description", "Content-Disposition"},
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

	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
	var swaggerSchemes []string
	schemes := os.Getenv("SWAGGER_SCHEME")

	if schemes == "" { // default to http and https
		swaggerSchemes = []string{"http", "https"}
	}

	docs.SwaggerInfo.Schemes = swaggerSchemes
	if os.Getenv("ENABLE_SWAGGER") == "true" {
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return r
}
