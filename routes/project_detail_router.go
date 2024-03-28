package routes

import (
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProjectDetailRouter(r *gin.RouterGroup) {
	projectDetailGroup := r.Group("/project_detail", middleware.RequiredAuth)
	{
		// projectDetailGroup.GET("/", controllers.GetProjectDetail)
		// projectDetailGroup.PUT("/", controllers.UpdateProjectDetail)
	}
}
