// routes/author_router.go

package routes

import (
	"vcs_backend/gorm/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAuthorRouter sets up the routes related to authors
func SetupAuthorRouter(r *gin.RouterGroup) {
	authorGroup := r.Group("/authors")
	{
		authorGroup.GET("/", controllers.AuthorController)
		// Add more author routes here if needed
	}
}
