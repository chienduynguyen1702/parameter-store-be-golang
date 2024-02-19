package initializers

import (
	"fmt"
	"vcs_backend/gorm/models"
)

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
