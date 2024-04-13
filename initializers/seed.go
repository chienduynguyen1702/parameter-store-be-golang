package initializers

import (
	"parameter-store-be/initializers/seed"

	"gorm.io/gorm"
)

func RunSeed(db *gorm.DB) error {
	// 1 - Seed organization , user, project, user_project_role, default role, stage, environment, permission
	if err := seed.SeedDatabase(db); err != nil {
		return err
	}

	// 2 - Seed data for test project
	if err := seed.SeedDataTestProjectUser(db); err != nil {
		return err
	}

	// 3 - Seed version, params for test project
	if err := seed.SeedDataTestProjectVersion(db); err != nil {
		return err
	}

	// 4 - Seed params for test project
	if err := seed.SeedDataTestParam(db); err != nil {
		return err
	}

	// 5 - Seed agent for test project
	if err := seed.SeedAgent(db); err != nil {
		return err
	}

	return nil
}
