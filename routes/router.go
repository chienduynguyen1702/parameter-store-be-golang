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

	// Define your routes here
	r.GET("/", controllers.MainController)
	r.GET("/posts", controllers.PostController)
	r.GET("/authors", controllers.AuthorController)

	return r
}
