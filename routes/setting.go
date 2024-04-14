package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupSetting(r *gin.RouterGroup) {
	userSettingGroup := r.Group("/settings", middleware.RequiredAuth)
	{
		userGroup := userSettingGroup.Group("/users")
		{
			userGroup.GET("", controllers.ListUser)
			userGroup.POST("", controllers.CreateUser)
			userGroup.GET("/:user_id", controllers.GetUserById)
			userGroup.PUT("/:user_id", controllers.UpdateUserInformation)
			userGroup.DELETE("/:user_id", controllers.DeleteUser)

			userGroup.GET("/archived", controllers.ListArchivedUser)
			userGroup.PATCH("/:user_id/archive", controllers.ArchiveUser)
			userGroup.PATCH("/:user_id/unarchive", controllers.RestoreUser)
		}
		roleGroup := userSettingGroup.Group("/roles")
		{
			roleGroup.GET("/", controllers.ListRole)
		}
	}
}
