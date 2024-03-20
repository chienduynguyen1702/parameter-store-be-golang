package initializers

import (
	"fmt"
	"log"
	"vcs_backend/gorm/models"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {

	err := db.AutoMigrate(&models.Organization{})
	if err != nil {
		log.Println("Failed to migrate Organization model")
		return err
	}
	err = db.AutoMigrate(&models.Project{})
	if err != nil {
		log.Println("Failed to migrate Project model")
		return err
	}
	err = db.AutoMigrate(&models.Version{})
	if err != nil {
		log.Println("Failed to migrate Version model")
		return err
	}
	err = db.AutoMigrate(&models.Stage{})
	if err != nil {
		log.Println("Failed to migrate Stage model")
		return err
	}
	err = db.AutoMigrate(&models.Environment{})
	if err != nil {
		log.Println("Failed to migrate Environment model")
		return err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Println("Failed to migrate User model")
		return err
	}
	err = db.AutoMigrate(&models.Token{})
	if err != nil {
		log.Println("Failed to migrate Token model")
		return err
	}
	err = db.AutoMigrate(&models.UserProjectRole{})
	if err != nil {
		log.Println("Failed to migrate UserProjectRole model")
		return err
	}
	err = db.AutoMigrate(&models.Role{})
	if err != nil {
		log.Println("Failed to migrate Role model")
		return err
	}
	err = db.AutoMigrate(&models.Permission{})
	if err != nil {
		log.Println("Failed to migrate Permission model")
		return err
	}
	err = db.AutoMigrate(&models.Agent{})
	if err != nil {
		log.Println("Failed to migrate Agent model")
		return err
	}
	err = db.AutoMigrate(&models.AgentLog{})
	if err != nil {
		log.Println("Failed to migrate AgentLog model")
	}
	fmt.Printf("\nDatabase migrated\n")
	return nil
}
