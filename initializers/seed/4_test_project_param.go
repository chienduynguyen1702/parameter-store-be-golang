package seed

import (
	"parameter-store-be/models"

	"gorm.io/gorm"
)

func SeedDataTestParam(db *gorm.DB) error {
	params := []models.Parameter{
		{
			Name:          "ENABLE_SWAGGER",
			Value:         "true",
			Description:   "To enable swagger",
			StageID:       7,
			EnvironmentID: 8,
			ProjectID:     4,
		},
		{
			Name:          "GIN_MODE",
			Value:         "release",
			Description:   "Set gin mode",
			StageID:       7,
			EnvironmentID: 8,
			ProjectID:     4,
		},
	}

	if err := db.Create(&params).Error; err != nil {
		return err
	}
	return nil
}
