package initializers

import (
	"fmt"
	"parameter-store-be/models"

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
		},
		{
			Name:        "Developer",
			Description: "Normal user",
		},
	}

	for _, role := range defaultRoles {
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}

	fmt.Printf("\nDatabase seeded\n")
	return nil
}
