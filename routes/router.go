// routes/router.go

package routes

import (
	"os"
	"parameter-store-be/controllers"
	docs "parameter-store-be/docs"
	"parameter-store-be/initializers"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up the routes for the application
func SetupRouter() *gin.Engine {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		SetupAuthorRouter(v1)
		SetupPostRouter(v1)

		// Add other route setups here if needed
	}

	docs.SwaggerInfo = initializers.SwaggerInfo
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1.GET("/helloworld", controllers.MainController)
	return r
}
