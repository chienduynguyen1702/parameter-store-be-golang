package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupStage(r *gin.RouterGroup) {
	stageGroup := r.Group("/stages", middleware.RequiredAuth)
	{
		stageGroup.GET("/", controllers.GetStages)
	}
}
