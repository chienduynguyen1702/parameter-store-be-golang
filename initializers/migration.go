package initializers

import (
	"fmt"
	"log"
	"parameter-store-be/models"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {

	err := db.AutoMigrate(&models.Organization{})
	if err != nil {
		log.Println("Failed to migrate Organization models")
		return err
	}
	err = db.AutoMigrate(&models.Project{})
	if err != nil {
		log.Println("Failed to migrate Project models")
		return err
	}
	err = db.AutoMigrate(&models.Version{})
	if err != nil {
		log.Println("Failed to migrate Version models")
		return err
	}
	err = db.AutoMigrate(&models.Stage{})
	if err != nil {
		log.Println("Failed to migrate Stage models")
		return err
	}
	err = db.AutoMigrate(&models.Environment{})
	if err != nil {
		log.Println("Failed to migrate Environment models")
		return err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Println("Failed to migrate User models")
		return err
	}
	err = db.AutoMigrate(&models.Token{})
	if err != nil {
		log.Println("Failed to migrate Token models")
		return err
	}
	err = db.AutoMigrate(&models.UserRoleProject{})
	if err != nil {
		log.Println("Failed to migrate UserRoleProject models")
		return err
	}
	err = db.AutoMigrate(&models.Role{})
	if err != nil {
		log.Println("Failed to migrate Role models")
		return err
	}
	err = db.AutoMigrate(&models.Permission{})
	if err != nil {
		log.Println("Failed to migrate Permission models")
		return err
	}
	err = db.AutoMigrate(&models.Agent{})
	if err != nil {
		log.Println("Failed to migrate Agent models")
		return err
	}
	err = db.AutoMigrate(&models.AgentLog{})
	if err != nil {
		log.Println("Failed to migrate AgentLog models")
	}
	fmt.Printf("\nDatabase migrated\n")
	return nil
}
