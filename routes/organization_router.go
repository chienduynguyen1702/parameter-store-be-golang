package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupOrganization(r *gin.RouterGroup) {
	organizationGroup := r.Group("/organization", middleware.RequiredAuth)
	{
		organizationGroup.GET("/", controllers.GetOrganizationInformation)
		organizationGroup.PUT("/", controllers.UpdateOrganizationInformation)
		organizationGroup.GET("/list-project", controllers.ListProjects)
		organizationGroup.POST("/new-project", controllers.CreateNewProject)
		organizationGroup.DELETE("/:project_id", controllers.DeleteProject)
	}
}
