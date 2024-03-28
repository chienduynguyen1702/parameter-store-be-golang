package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

// SetupProjectRouter sets up the routes related to authors
func SetupProjectRouter(r *gin.RouterGroup) {
	projectGroup := r.Group("/projects", middleware.RequiredAuth)
	{
		projectGroup.GET("/", controllers.ListProjects)
		projectGroup.POST("/", controllers.CreateNewProject)

		// projectGroup.GET("/:project_id", controllers.GetProjectInformation)
		projectGroup.PUT("/:project_id", controllers.UpdateProjectInformation)
		projectGroup.DELETE("/:project_id", controllers.DeleteProject)
	}
}
