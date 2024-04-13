package seed

import (
	"log"
	"parameter-store-be/models"
	"time"

	"gorm.io/gorm"
)

// required organization is seeded in the first migration

func SeedDataTestProjectUser(db *gorm.DB) error {
	testestUser := models.User{
		Username:            "test user",
		Email:               "test@gmail.com",
		Password:            "$2a$10$qdi5VjamNQsbgisE7ijEx.McxvM5eQzCcDmvDosm5cSDhwkznMOCa", // 123123
		OrganizationID:      sampleOrganizations.ID,
		IsOrganizationAdmin: false,
		Phone:               "0123123123",
	}
	if err := db.Create(&testestUser).Error; err != nil {
		return err
	}
	testUser = testestUser
	log.Printf("\nTest user data is seeded.\n")

	golang_swagger := models.Project{
		Name:           "Golang Swagger Project",
		StartAt:        time.Now(),
		Description:    "Test project description",
		CurrentSprint:  "1",
		Status:         "In Progress",
		RepoURL:        "github.com/chienduynguyen1702/golang-swagger",
		OrganizationID: sampleOrganizations.ID,
		Address:        "SoICT, HUST",
		RepoApiToken:   "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",
	}

	if err := db.Create(&golang_swagger).Error; err != nil {
		return err
	}
	testProject = golang_swagger
	upr := []models.UserRoleProject{
		// admin user as project admin
		{
			UserID:    sampleAdmin.ID,
			ProjectID: testProject.ID,
			RoleID:    defaultRoles[1].ID,
		},
		// test user as project member
		{
			UserID:    testUser.ID,
			ProjectID: testProject.ID,
			RoleID:    defaultRoles[2].ID,
		},
	}
	if err := db.Create(&upr).Error; err != nil {
		return err
	}
	log.Printf("\nTest project data is seeded.\n")
	return nil
}
