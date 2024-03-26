package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

// SetupProjectRouter sets up the routes related to authors
func SetupProjectRouter(r *gin.RouterGroup) {
	projectGroup := r.Group("/projects")
	{
		projectGroup.GET("/", middleware.RequiredAuth, controllers.ListProjects)
		projectGroup.POST("/", middleware.RequiredAuth, controllers.CreateNewProject)

		// projectGroup.GET("/:project_id", middleware.RequiredAuth, controllers.GetProjectInformation)
		projectGroup.PUT("/:project_id", middleware.RequiredAuth, controllers.UpdateProjectInformation)
		projectGroup.DELETE("/:project_id", middleware.RequiredAuth, controllers.DeleteProject)
	}
}
