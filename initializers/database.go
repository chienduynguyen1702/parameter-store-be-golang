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
func SeedDatabase() {
	authors := []*models.Author{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@gmail.com",
			Password:  "password",
			Phone:     "1234567890",
			Address:   "123 Main Street",
			Posts: []models.Post{
				{
					Title: "First Post",
					Body:  "This is the first post",
				},
				{
					Title: "Second Post",
					Body:  "This is the second post",
				},
			},
		},
		{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "janedoegmail.com",
			Password:  "password",
			Phone:     "1231234123490",
			Address:   "123 Main Street",
			Posts: []models.Post{
				{
					Title: "Third Post",
					Body:  "This is the third post",
				},
			},
		},
		{
			FirstName: "John",
			LastName:  "Smith",
			Email:     "johnsmith@gmail.com",
			Phone:     "1234567890",
			Password:  "password",
			Address:   "456 Main Street",
			Posts:     []models.Post{},
		},
	}
	DB.Create(authors)

	fmt.Printf("\nDatabase seeded\n")
}

func QueryPosts() []models.Post {
	var posts []models.Post
	DB.Find(&posts)
	return posts
}
