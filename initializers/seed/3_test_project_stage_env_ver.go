package seed

import (
	"log"
	"parameter-store-be/models"

	"gorm.io/gorm"
)

func SeedDataTestProjectStageEnvVer(db *gorm.DB) error {
	stages := []models.Stage{
		{
			Name:        "Build",
			Description: "Build stage",
			// ProjectID:   testProject.ID,
			ProjectID: 4,
		},
		{
			Name:        "Deploy",
			Description: "Deploy stage",
			// ProjectID:   testProject.ID,
			ProjectID: 4,
		},
	}
	if err := db.Create(&stages).Error; err != nil {
		return err
	}
	log.Printf("\nTest project stages data is seeded.\n")

	envs := []models.Environment{
		{
			Name:        "Development",
			Description: "Development environment",
			// ProjectID:   testProject.ID,
			ProjectID: 4,
		},
		{
			Name:        "Production",
			Description: "Production environment",
			// ProjectID:   testProject.ID,
			ProjectID: 4,
		},
	}
	if err := db.Create(&envs).Error; err != nil {
		return err
	}
	log.Printf("\nTest project environments data is seeded.\n")

	vers := []models.Version{
		{
			Name:        "v1.0",
			Description: "Version 1.0",
			// ProjectID:   testProject.ID,
			ProjectID: 4,
		},
		{
			Name:        "v1.5",
			Description: "Version 1.5",
			// ProjectID:   testProject.ID,
			ProjectID: 4,
		},
	}
	if err := db.Create(&vers).Error; err != nil {
		return err
	}

	return nil
}
