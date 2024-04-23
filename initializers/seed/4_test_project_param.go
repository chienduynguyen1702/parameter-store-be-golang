package seed

import (
	"log"
	"parameter-store-be/models"

	"gorm.io/gorm"
)

func SeedDataTestParam(db *gorm.DB) error {
	params := []models.Parameter{
		{
			Name:          "ENABLE_SWAGGER",
			Value:         "true",
			Description:   "To enable swagger",
			StageID:       defaultStages[3].ID,
			EnvironmentID: defaultEnvironments[3].ID,
			ProjectID:     testProject.ID,
		},
		{
			Name:          "GIN_MODE",
			Value:         "release",
			Description:   "Set gin mode",
			StageID:       defaultStages[3].ID,
			EnvironmentID: defaultEnvironments[3].ID,
			ProjectID:     testProject.ID,
		},
	}

	if err := db.Create(&params).Error; err != nil {
		return err
	}
	testParams = params

	testVersions[0].Parameters = append(testVersions[0].Parameters, testParams[0])
	testVersions[1].Parameters = append(testVersions[1].Parameters, testParams[0], testParams[1])

	if err := db.Save(&testVersions[0]).Error; err != nil {
		return err
	}
	if err := db.Save(&testVersions[1]).Error; err != nil {
		return err
	}
	log.Println("Test project params are seeded.")
	return nil
}
