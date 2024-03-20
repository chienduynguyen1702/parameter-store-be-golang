// routes/router.go

package routes

import (
	"parameter-store-be/controllers"
	docs "parameter-store-be/docs"
	"parameter-store-be/initializers"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up the routes for the application
func SetupRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/api/v1", controllers.MainController)
	v1 := r.Group("/api/v1")
	{
		SetupAuthorRouter(v1)
		SetupPostRouter(v1)

		// Add other route setups here if needed
	}

	docs.SwaggerInfo = initializers.SwaggerInfo
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
