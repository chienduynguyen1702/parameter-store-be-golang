// routes/author_router.go

package routes

import (
	pc "vcs_backend/gorm/controllers/post"

	"github.com/gin-gonic/gin"
)

// SetupAuthorRouter sets up the routes related to authors
func SetupPostRouter(r *gin.RouterGroup) {
	postGroup := r.Group("/posts")
	{
		postGroup.GET("/", pc.GetPosts)
		postGroup.GET("/:id", pc.GetPostByID)
		// postGroup.GET("/", pc.GetPostsByAuthorID)
		// postGroup.POST("/", pc.CreateNewPost)
	}
}
