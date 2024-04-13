package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupEnvironment(r *gin.RouterGroup) {
	environmentGroup := r.Group("/envs", middleware.RequiredAuth)
	{
		environmentGroup.GET("/", controllers.GetEnvironments)
	}
}
