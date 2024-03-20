package main

import (
	"fmt"
	"log"
	"os"
	"vcs_backend/gorm/controllers"
	"vcs_backend/gorm/initializers"
	"vcs_backend/gorm/routes"
)

func init() {
	initializers.LoadEnvVariables()
	db, err := initializers.ConnectDatabase() // return *gorm.DB
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	initializers.Migration(db) // migration db
	controllers.SetDB(db)      // set controller use that db *gorm.DB
	// initializers.SeedDatabase()
	fmt.Println("")

}
func main() {
	r := routes.SetupRouter()
	fmt.Println("Server is running on port", os.Getenv("PORT"))
	r.Run(":" + os.Getenv("PORT"))
}
