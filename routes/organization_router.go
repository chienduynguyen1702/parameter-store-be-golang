package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupOrganization(r *gin.RouterGroup) {
	organizationGroup := r.Group("/organizations", middleware.RequiredAuth)
	{
		organizationGroup.GET("/", controllers.GetOrganizationInformation)
		organizationGroup.PUT("/", controllers.UpdateOrganizationInformation, middleware.RequiredIsOrgAdmin)
		organizationGroup.GET("/projects", controllers.ListProjects)
		organizationGroup.POST("/projects", controllers.CreateNewProject, middleware.RequiredIsOrgAdmin)
		organizationGroup.DELETE("/projects/:project_id", controllers.DeleteProject, middleware.RequiredIsOrgAdmin)
	}
}
