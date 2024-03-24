package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

// SetupOrganizationRouter sets up the routes related to authors
func SetupOrganizationRouter(r *gin.RouterGroup) {
	authGroup := r.Group("/organization")
	{
		authGroup.GET("/:organization_id", middleware.RequiredAuth, controllers.GetOrganization)
		// authGroup.GET("/:organization_id", controllers.GetOrganization)
		// authGroup.POST("/", controllers.CreateNewPost)
	}
}
