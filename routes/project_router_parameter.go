package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupProjectParameter(r *gin.RouterGroup) {
	projectGroup := r.Group("/project/param", middleware.RequiredAuth)
	{
		// projectGroup.GET("/:project_id", controllers.GetProjectInformation)
		projectGroup.PUT("/:project_id", controllers.UpdateProjectInformation)
	}
}
