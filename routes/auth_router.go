package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthorRouter sets up the routes related to authors
func SetupAuthRouter(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", controllers.Login)
		authGroup.POST("/register", controllers.Register)
		authGroup.GET("/validate", middleware.RequiredAuth, controllers.Validate)
		// authGroup.POST("/", controllers.CreateNewPost)
	}
}
