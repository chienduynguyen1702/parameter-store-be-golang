package seed

import (
	"parameter-store-be/models"

	"gorm.io/gorm"
)

func SeedAgent(DB *gorm.DB) error {
	agent := models.Agent{
		ProjectID:     4,
		Name:          "Test Agent",
		APIToken:      "123123", // use this token to Get params from agent by stage and environment
		StageID:       8,
		EnvironmentID: 8,
		WorkflowName:  "Build Docker And Deploy",
	}
	if err := DB.Create(&agent).Error; err != nil {
		return err
	}
	return nil
}
