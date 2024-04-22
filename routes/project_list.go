package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupProjectList(r *gin.RouterGroup) {
	projectListGroup := r.Group("/project-list", middleware.RequiredAuth)
	{
		projectListGroup.GET("/", controllers.ListProjects)
		projectListGroup.POST("/", middleware.RequiredIsOrgAdmin, controllers.CreateNewProject)
		// projectListGroup.DELETE("/:project_id", middleware.RequiredIsOrgAdmin, controllers.DeleteProject)

		projectListGroup.GET("/archived", controllers.ListArchivedProjects)
		projectListGroup.PATCH("/:project_id/archive", middleware.RequiredIsOrgAdmin, controllers.ArchiveProject)
		projectListGroup.PATCH("/:project_id/unarchive", middleware.RequiredIsOrgAdmin, controllers.UnarchiveProject)
	}
}
