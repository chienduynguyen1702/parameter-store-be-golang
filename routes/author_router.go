// routes/author_router.go

package routes

import (
	ac "vcs_backend/gorm/controllers/author"

	"github.com/gin-gonic/gin"
)

// SetupAuthorRouter sets up the routes related to authors
func SetupAuthorRouter(r *gin.RouterGroup) {
	authorGroup := r.Group("/authors")
	{
		// CREATE
		authorGroup.POST("/register", ac.RegisterAuthor)
		authorGroup.POST("/:id/publish", ac.CreateNewPost)
		// READ
		authorGroup.GET("/by-id", ac.GetAuthorById)
		authorGroup.GET("/by-name", ac.GetAuthorsByName)
		authorGroup.GET("/", ac.GetAllAuthors)
		// UPDATE
		authorGroup.PUT("/:id", ac.UpdateAuthorInfo)

		// DELETE
		authorGroup.DELETE("/:id", ac.DeleteAuthor)
	}
}
