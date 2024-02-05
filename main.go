package main

import (
	"vcs_backend/gorm/controllers"
	"vcs_backend/gorm/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()

}
func main() {
	r := gin.Default()
	// gin.SetMode(gin.ReleaseMode)
	// routes.UserRoute(router)
	r.GET("/post", controllers.Posts)
	// fmt.Println("Server is running on port ", os.Getenv("PORT"))
	r.Run()
}
