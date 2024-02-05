package routes

import (
	uc "vcs_backend/gorm/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.GET("/users", uc.UserController)
}

// router.GET("/users/:id", controllers.GetUser)
// router.POST("/users", controllers.CreateUser)
// router.PUT("/users/:id", controllers.UpdateUser)
// router.DELETE("/users/:id", controllers.DeleteUser)
