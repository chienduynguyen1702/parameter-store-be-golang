package initializers

import (
	"parameter-store-be/initializers/seed"

	"gorm.io/gorm"
)

func RunSeed(db *gorm.DB) error {
	// 1 - Seed organization , user, project, user_project_role
	if err := seed.SeedDatabase(db); err != nil {
		return err
	}

	// 2 - Seed data for test project
	if err := seed.SeedDataTestProjectUser(db); err != nil {
		return err
	}

	// 3 - Seed stages, environment, version, params for test project
	if err := seed.SeedDataTestProjectStageEnvVer(db); err != nil {
		return err
	}

	return nil
}
