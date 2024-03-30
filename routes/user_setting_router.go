package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupSetting(r *gin.RouterGroup) {
	userSettingGroup := r.Group("/setting", middleware.RequiredAuth)
	{
		userGroup := userSettingGroup.Group("/user")
		{
			userGroup.GET("/", controllers.ListUser)
			userGroup.POST("/", controllers.CreateUser)
			userGroup.PUT("/:user_id", controllers.UpdateUserInformation)
			userGroup.DELETE("/:user_id", controllers.DeleteUser)
		}
		roleGroup := userSettingGroup.Group("/role")
		{
			roleGroup.GET("/", controllers.ListRole)
		}
	}
}
