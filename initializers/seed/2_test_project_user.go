package seed

import (
	"log"
	"parameter-store-be/models"
	"time"

	"gorm.io/gorm"
)

// required organization is seeded in the first migration

func SeedDataTestProjectUser(db *gorm.DB) error {
	testUser := models.User{
		Username:            "test user",
		Email:               "test@gmail.com",
		Password:            "$2a$10$KCzNv5lThy0h65JVRp/.huq1kxa6oq5jt.OHyqy6YfpBd4TAKIk3C",
		OrganizationID:      1,
		IsOrganizationAdmin: false,
		Phone:               "0123123123",
	}
	if err := db.Create(&testUser).Error; err != nil {
		return err
	}
	log.Printf("\nTest user data is seeded.\n")

	testProject := models.Project{
		Name:           "Test Project",
		StartAt:        time.Now(),
		Description:    "Test project description",
		CurrentSprint:  "1",
		Status:         "In Progress",
		RepoURL:        "github.com/chienduynguyen1702/golang-swagger",
		OrganizationID: 1,
		Address:        "SoICT, HUST",
		RepoApiToken:   "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",
	}

	if err := db.Create(&testProject).Error; err != nil {
		return err
	}
	upr := []models.UserRoleProject{
		// admin user as project admin
		{
			UserID:    1,
			ProjectID: testProject.ID,
			RoleID:    2,
		},
		// test user as project member
		{
			UserID:    testUser.ID,
			ProjectID: testProject.ID,
			RoleID:    3,
		},
	}
	if err := db.Create(&upr).Error; err != nil {
		return err
	}
	log.Printf("\nTest project data is seeded.\n")
	return nil
}
