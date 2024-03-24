package initializers

import (
	"fmt"
	"parameter-store-be/models"

	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	defaultRoles := []models.Role{
		{
			Name: "",
		},
		{
			Name: "user",
		},
	}

	for _, role := range defaultRoles {
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}

	fmt.Printf("\nDatabase seeded\n")
	return nil
}
