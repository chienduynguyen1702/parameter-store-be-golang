package main

import (
	"log"
	"os"
	"parameter-store-be/controllers"
	"parameter-store-be/initializers"
	"parameter-store-be/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	if os.Getenv("SERVERLESS_DEPLOY") != "true" {
		initializers.LoadEnvVariables()
	}
	db, err := initializers.ConnectDatabase() // return *gorm.DB
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Migration db
	if os.Getenv("RUN_MIGRATION") == "true" {
		initializers.Migration(db) // migration db
	}
	// Seed data
	if os.Getenv("RUN_SEED") == "true" {
		if err := initializers.RunSeed(db); err != nil {
			log.Fatal("Failed to seed database")
		}
	}

	// Set controller
	controllers.SetDB(db) // set controller use that db *gorm.DB
	log.Println("Finished init.")
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
	if os.Getenv("ENABLE_HEALTH_CHECK") == "true" {
		go controllers.ScheduleWorkflowCheck()
	}
	log.Println("Server is running on port", port, "in", os.Getenv("GIN_MODE"), "mode.")
	r.Run(":" + port)
}
