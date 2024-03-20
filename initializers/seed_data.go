package initializers

import (
	"fmt"

	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {

	fmt.Printf("\nDatabase seeded\n")
	return nil
}
