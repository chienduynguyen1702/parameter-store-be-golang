package seed

import (
	"log"
	"parameter-store-be/models"

	"gorm.io/gorm"
)

func SeedAgent(DB *gorm.DB) error {
	agent := models.Agent{
		ProjectID:     testProject.ID,
		Name:          "Test Agent",
		APIToken:      "123123", // use this token to Get params from agent by stage and environment
		StageID:       defaultStages[3].ID,
		EnvironmentID: defaultEnvironments[3].ID,
		WorkflowName:  "Build Docker And Deploy",
	}
	if err := DB.Create(&agent).Error; err != nil {
		return err
	}
	log.Println("Test project agent is seeded.")
	return nil
}
