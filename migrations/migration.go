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
	err = initializers.DB.AutoMigrate(&models.Author{})
	if err != nil {
		panic("Failed to migrate Author model")
	}
	err = initializers.DB.AutoMigrate(&models.Author_Post{})
	if err != nil {
		panic("Failed to migrate Author_Post model")
	}

}
