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
			userGroup.POST("", middleware.RequiredIsAdmin, controllers.CreateUser)
			userGroup.GET("/:user_id", controllers.GetUserById)
			userGroup.PUT("/:user_id", middleware.RequiredIsAdmin, controllers.UpdateUserInformation)
			userGroup.DELETE("/:user_id", middleware.RequiredIsAdmin, controllers.DeleteUser)

			userGroup.GET("/archived", controllers.ListArchivedUser)
			userGroup.PATCH("/:user_id/archive", middleware.RequiredIsAdmin, controllers.ArchiveUser)
			userGroup.PATCH("/:user_id/unarchive", middleware.RequiredIsAdmin, controllers.RestoreUser)
		}
		roleGroup := userSettingGroup.Group("/roles")
		{
			roleGroup.GET("/", controllers.ListRole)
		}
	}
}
