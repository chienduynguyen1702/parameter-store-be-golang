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
	db := initializers.ConnectDatabase() // return *gorm.DB
	controllers.SetDB(db)                // set controller use that db *gorm.DB
	// initializers.SeedDatabase()
	fmt.Println("")

}
func main() {
	r := routes.SetupRouter()
	fmt.Println("Server is running on port", os.Getenv("PORT"))
	r.Run(":" + os.Getenv("PORT"))
}
