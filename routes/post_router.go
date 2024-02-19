// routes/author_router.go

package routes

import (
	"vcs_backend/gorm/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAuthorRouter sets up the routes related to authors
func SetupPostRouter(r *gin.RouterGroup) {
	postGroup := r.Group("/posts")
	{
		postGroup.GET("/", controllers.PostController)
		// Add more author routes here if needed
	}
}
