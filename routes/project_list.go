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
		projectListGroup.POST("/", controllers.CreateNewProject, middleware.RequiredIsOrgAdmin)
		// projectListGroup.DELETE("/:project_id", controllers.DeleteProject, middleware.RequiredIsOrgAdmin)

		projectListGroup.GET("/archived", controllers.ListArchivedProjects)
		projectListGroup.PATCH("/:project_id/archive", controllers.ArchiveProject, middleware.RequiredIsOrgAdmin)
		projectListGroup.PATCH("/:project_id/unarchive", controllers.UnarchiveProject, middleware.RequiredIsOrgAdmin)
	}
}
