// routes/author_router.go

package routes

import (
	"parameter-store-be/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAuthorRouter sets up the routes related to authors
func SetupPostRouter(r *gin.RouterGroup) {
	postGroup := r.Group("/posts")
	{
		postGroup.GET("/", controllers.GetPosts)
		postGroup.GET("/by-author-id/:author-id", controllers.GetPostsByAuthorID)
		postGroup.GET("/:id", controllers.GetPostByID)
		// postGroup.GET("/", controllers.GetPostsByAuthorID)
		// postGroup.POST("/", controllers.CreateNewPost)
	}
}
