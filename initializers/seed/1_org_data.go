package seed

import (
	"log"
	"parameter-store-be/models"
	"time"

	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	allPermissions := []models.Permission{

		{
			Name:        "auth",
			Description: "Authentication",
		},
		{
			Name:        "user-update",
			Description: "Update user",
		},
		{
			Name:        "user-update-avatar",
			Description: "Upload image",
		},
		{
			Name:        "user-create",
			Description: "Create user",
		},
		{
			Name:        "user-one",
			Description: "Get user",
		},
		{
			Name:        "user-list",
			Description: "Get users",
		},
		{
			Name:        "user-archivist-archive",
			Description: "Archive users",
		},
		{
			Name:        "user-archivist-unarchive",
			Description: "Unarchive users",
		},
		{
			Name:        "user-archivist-list",
			Description: "Get archived users",
		},
		{
			Name:        "role-create",
			Description: "Create role",
		},
		{
			Name:        "role-update",
			Description: "Update role",
		},
		{
			Name:        "role-one",
			Description: "Get role",
		},
		{
			Name:        "role-list",
			Description: "Get roles",
		},
		{
			Name:        "role-archivist-archive",
			Description: "Archive roles",
		},
		{
			Name:        "role-archivist-unarchive",
			Description: "Unarchive roles",
		},
		{
			Name:        "role-archivist-list",
			Description: "Get archived roles",
		},
		{
			Name:        "permission-list",
			Description: "Get permissions",
		},
		{
			Name:        "setting-list",
			Description: "Get settings",
		},
		{
			Name:        "setting-update",
			Description: "Update settings",
		},
		{
			Name:        "setting-create",
			Description: "Create settings",
		},
		{
			Name:        "content-list",
			Description: "Get contents",
		},
		{
			Name:        "user-summary",
			Description: "Get user summary",
		},
		{
			Name:        "user-export",
			Description: "Export users",
		},
	}

	for _, permission := range allPermissions {
		if err := db.Create(&permission).Error; err != nil {
			return err
		}
	}
	defaultRoles := []models.Role{
		{
			Name:        "Organization Admin",
			Description: "Admin of the organization",
			Permissions: allPermissions,
		},
		{
			Name:        "Project Admin",
			Description: "Admin of the project",
			Permissions: allPermissions,
		},
		{
			Name:        "Developer",
			Description: "Normal user",
			Permissions: allPermissions,
		},
	}

	for i, role := range defaultRoles {
		if err := db.Create(&role).Error; err != nil {
			return err
		}
		defaultRoles[i] = role
	}

	log.Printf("\nDefault roles and permission data are seeded.\n")

	organization := models.Organization{
		Name:              "HUST",
		AliasName:         "Hanoi University of Science and Technology",
		Address:           "1 Dai Co Viet, Hanoi",
		EstablishmentDate: time.Date(1956, time.Month(10), 10, 0, 0, 0, 0, time.UTC),
		Description:       "Hanoi University of Science and Technology (HUST) is a multidisciplinary technical university located in Hanoi, Vietnam. It was founded on October 10, 1956, and is one of the two largest technical universities in Vietnam.",
	}

	if err := db.Create(&organization).Error; err != nil {
		return err
	}

	log.Printf("\nDefault organization data is seeded.\n")

	admin := models.User{
		Username:            "admin",
		Email:               "admin@gmail.com",
		Password:            "$2a$10$qdi5VjamNQsbgisE7ijEx.McxvM5eQzCcDmvDosm5cSDhwkznMOCa", // 123123
		OrganizationID:      organization.ID,
		IsOrganizationAdmin: true,
		Phone:               "0123456789",
		// LastLogIn:           time.Now(),
	}
	user1 := models.User{
		Username:            "user1",
		Email:               "user1@gmail.com",
		Password:            "$2a$10$qdi5VjamNQsbgisE7ijEx.McxvM5eQzCcDmvDosm5cSDhwkznMOCa", // 123123
		OrganizationID:      organization.ID,
		IsOrganizationAdmin: false,
		Phone:               "0123451231",
		// LastLogIn:           time.Now(),
	}
	user2 := models.User{
		Username:            "user2",
		Email:               "user2@gmail.com",
		Password:            "$2a$10$qdi5VjamNQsbgisE7ijEx.McxvM5eQzCcDmvDosm5cSDhwkznMOCa", // 123123
		OrganizationID:      organization.ID,
		IsOrganizationAdmin: false,
		Phone:               "0123451232",
		// LastLogIn:           time.Now(),
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	if err := db.Create(&user1).Error; err != nil {
		return err
	}
	if err := db.Create(&user2).Error; err != nil {
		return err
	}
	log.Printf("\nDefault user data is seeded.\n")

	projects := []models.Project{
		{
			Name:           "Parameter Store",
			StartAt:        time.Now(),
			Description:    "Parameter Store is a project to store parameters",
			CurrentSprint:  "1",
			Status:         "In Progress",
			RepoURL:        "github.com/parameter-store",
			OrganizationID: organization.ID,
			Address:        "SoICT, HUST",
		},
		{
			Name:           "Parameter Store 2",
			StartAt:        time.Now(),
			Description:    "Parameter Store is a project to store parameters",
			CurrentSprint:  "1",
			Status:         "In Progress",
			RepoURL:        "github.com/parameter-store",
			OrganizationID: organization.ID,
			Address:        "SoICT, HUST",
		},
		{
			Name:           "Parameter Store 3",
			StartAt:        time.Now(),
			Description:    "Parameter Store is a project to store parameters",
			CurrentSprint:  "1",
			Status:         "In Progress",
			RepoURL:        "github.com/parameter-store",
			OrganizationID: organization.ID,
			Address:        "SoICT, HUST",
		},
	}

	for i, project := range projects {
		if err := db.Create(&project).Error; err != nil {
			return err
		}
		projects[i] = project
	}
	log.Printf("\nDefault project data is seeded.\n")

	upr := []models.UserRoleProject{
		// project 1
		{
			UserID:    admin.ID,
			ProjectID: projects[0].ID,
			RoleID:    defaultRoles[1].ID,
		},
		{
			UserID:    user1.ID,
			ProjectID: projects[0].ID,
			RoleID:    defaultRoles[2].ID,
		},

		// project 2
		{
			UserID:    admin.ID,
			ProjectID: projects[1].ID,
			RoleID:    defaultRoles[1].ID,
		},
		{
			UserID:    user1.ID,
			ProjectID: projects[1].ID,
			RoleID:    defaultRoles[2].ID,
		},

		// project 3
		{
			UserID:    admin.ID,
			ProjectID: projects[2].ID,
			RoleID:    defaultRoles[1].ID,
		},
		{
			UserID:    user2.ID,
			ProjectID: projects[2].ID,
			RoleID:    defaultRoles[2].ID,
		},
	}

	for _, upr := range upr {
		if err := db.Create(&upr).Error; err != nil {
			return err
		}
	}

	log.Printf("\nDefault relation user project role is seed\n")

	stages := []models.Stage{
		{
			Name:        "Build",
			Description: "Build stage",
			ProjectID:   projects[0].ID,
		},
		{
			Name:        "Deploy",
			Description: "Deploy stage",
			ProjectID:   projects[0].ID,
		},
		{
			Name:        "Build",
			Description: "Build stage",
			ProjectID:   projects[1].ID,
		},
		{
			Name:        "Deploy",
			Description: "Deploy stage",
			ProjectID:   projects[1].ID,
		},
		{
			Name:        "Build",
			Description: "Build stage",
			ProjectID:   projects[2].ID,
		},
		{
			Name:        "Deploy",
			Description: "Deploy stage",
			ProjectID:   projects[2].ID,
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
			ProjectID:   projects[0].ID,
		},
		{
			Name:        "Production",
			Description: "Production environment",
			ProjectID:   projects[0].ID,
		},
		{
			Name:        "Development",
			Description: "Development environment",
			ProjectID:   projects[1].ID,
		},
		{
			Name:        "Production",
			Description: "Production environment",
			ProjectID:   projects[1].ID,
		},
		{
			Name:        "Development",
			Description: "Development environment",
			ProjectID:   projects[2].ID,
		},
		{
			Name:        "Production",
			Description: "Production environment",
			ProjectID:   projects[2].ID,
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
			ProjectID:   projects[0].ID,
		},
		{
			Name:        "v1.5",
			Description: "Version 1.5",
			ProjectID:   projects[0].ID,
		},
		{
			Name:        "v1.0",
			Description: "Version 1.0",
			ProjectID:   projects[1].ID,
		},
		{
			Name:        "v1.5",
			Description: "Version 1.5",
			ProjectID:   projects[1].ID,
		},
		{
			Name:        "v1.0",
			Description: "Version 1.0",
			ProjectID:   projects[2].ID,
		},
		{
			Name:        "v1.5",
			Description: "Version 1.5",
			ProjectID:   projects[2].ID,
		},
	}
	if err := db.Create(&vers).Error; err != nil {
		return err
	}
	return nil
}
