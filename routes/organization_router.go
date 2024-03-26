package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

// SetupOrganizationRouter sets up the routes related to authors
func SetupOrganizationRouter(r *gin.RouterGroup) {
	organizationGroup := r.Group("/organization")
	{
		organizationGroup.GET("/", middleware.RequiredAuth, controllers.GetOrganizationInformation)
		organizationGroup.PUT("/", middleware.RequiredAuth, controllers.UpdateOrganizationInformation)
		
	}
}
