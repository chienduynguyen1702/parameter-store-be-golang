package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

// SetupOrganizationRouter sets up the routes related to authors
func SetupOrganizationRouter(r *gin.RouterGroup) {
	organizationGroup := r.Group("/organization", middleware.RequiredAuth)
	{
		organizationGroup.GET("/", controllers.GetOrganizationInformation)
		organizationGroup.PUT("/", controllers.UpdateOrganizationInformation)

	}
}
