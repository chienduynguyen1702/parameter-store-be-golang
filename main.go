package main

import (
	"fmt"
	"os"
	"vcs_backend/gorm/controllers"
	"vcs_backend/gorm/initializers"
	"vcs_backend/gorm/routes"
)

func init() {
	initializers.LoadEnvVariables()
	db := initializers.ConnectDatabase() // Modify this line
	controllers.SetDB(db)                // Add this line
	// initializers.SeedDatabase()
	fmt.Println("")

}
func main() {
	r := routes.SetupRouter()
	fmt.Println("Server is running on port", os.Getenv("PORT"))
	r.Run(":" + os.Getenv("PORT"))
}
