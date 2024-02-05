package main

import (
	"vcs_backend/gorm/initializers"
	"vcs_backend/gorm/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
}
func main() {
	err := initializers.DB.AutoMigrate(&models.Post{})
	if err != nil {
		panic("Failed to migrate Post model")
	}

}
