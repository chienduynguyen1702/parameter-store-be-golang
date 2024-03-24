package main

import (
	"fmt"
	"log"
	"os"
	"parameter-store-be/controllers"
	"parameter-store-be/initializers"
	"parameter-store-be/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	db, err := initializers.ConnectDatabase() // return *gorm.DB
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	initializers.Migration(db) // migration db
	controllers.SetDB(db)      // set controller use that db *gorm.DB
	// if err := initializers.SeedDatabase(db); err != nil {
	// 	log.Fatal("Failed to seed database")
	// }
}
func main() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := routes.SetupV1Router()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	fmt.Println("Server is running on port", port)
	r.Run(":" + port)
}
