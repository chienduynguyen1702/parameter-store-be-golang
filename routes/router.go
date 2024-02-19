// routes/router.go

package routes

import (
	"vcs_backend/gorm/controllers"

	"github.com/gin-gonic/gin"
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

	return r
}
