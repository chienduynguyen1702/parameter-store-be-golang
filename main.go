package main

import (
	"fmt"
	"os"
	"vcs_backend/gorm/controllers"
	"vcs_backend/gorm/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	db := initializers.ConnectDatabase() // Modify this line
	controllers.SetDB(db)                // Add this line
	// initializers.SeedDatabase()

}
func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	// routes.UserRoute(router)
	r.GET("/posts", controllers.PostController)
	fmt.Println("Server is running on port ", os.Getenv("PORT"))
	r.Run()
}
