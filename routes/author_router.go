// routes/author_router.go

package routes

import (
	"parameter-store-be/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAuthorRouter sets up the routes related to authors
func SetupAuthorRouter(r *gin.RouterGroup) {
	authorGroup := r.Group("/authors")
	{
		// CREATE
		authorGroup.POST("/register", controllers.RegisterAuthor)
		authorGroup.POST("/:id/publish", controllers.CreateNewPost)
		// READ
		authorGroup.GET("/by-id", controllers.GetAuthorById)
		authorGroup.GET("/by-name", controllers.GetAuthorsByName)
		authorGroup.GET("/", controllers.GetAllAuthors)
		// UPDATE
		authorGroup.PUT("/:id", controllers.UpdateAuthorInfo)

		// DELETE
		authorGroup.DELETE("/:id", controllers.DeleteAuthor)
		authorGroup.DELETE("/delete-post/:author-id", controllers.DeletePostOfAuthor)
	}
}
