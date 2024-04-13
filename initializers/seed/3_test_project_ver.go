package seed

import (
	"parameter-store-be/models"

	"gorm.io/gorm"
)

func SeedDataTestProjectVersion(db *gorm.DB) error {

	vers := []models.Version{
		{
			Name:        "v1.0",
			Description: "Version 1.0",
			ProjectID:   testProject.ID,
		},
		{
			Name:        "v1.5",
			Description: "Version 1.5",
			ProjectID:   testProject.ID,
		},
	}
	if err := db.Create(&vers).Error; err != nil {
		return err
	}
	testVersions = vers
	return nil
}
