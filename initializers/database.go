package initializers

import (
	"fmt"
	"os"
	"vcs_backend/gorm/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
	)
	// dsn := "host=localhost user=postgres password=postgres dbname=learning_gorm port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Printf("\nDatabase connected\n")
	err = DB.AutoMigrate(&models.Author{})
	if err != nil {
		panic("Failed to migrate Author model")
	}
	err = DB.AutoMigrate(&models.Post{})
	if err != nil {
		panic("Failed to migrate Post model")
	}
	fmt.Printf("\nDatabase migrated\n")
	// DB.AutoMigrate(&models.Author_Post{})
	return DB
}