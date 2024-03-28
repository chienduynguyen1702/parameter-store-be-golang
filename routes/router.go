// routes/router.go

package routes

import (
	docs "parameter-store-be/docs"
	"parameter-store-be/initializers"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath  /api/v1
// SetupRouter sets up the routes for the application
func SetupV1Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo = initializers.SwaggerInfo
	v1 := r.Group("/api/v1")
	{
		SetupAuthRouter(v1)
		SetupOrganizationRouter(v1)
		SetupProjectRouter(v1)
		SetupProjectDetailRouter(v1)
	}

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
