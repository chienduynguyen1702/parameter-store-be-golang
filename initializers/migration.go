package initializers

import (
	"fmt"
	"log"
	"parameter-store-be/model"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {

	err := db.AutoMigrate(&model.Organization{})
	if err != nil {
		log.Println("Failed to migrate Organization model")
		return err
	}
	err = db.AutoMigrate(&model.Project{})
	if err != nil {
		log.Println("Failed to migrate Project model")
		return err
	}
	err = db.AutoMigrate(&model.Version{})
	if err != nil {
		log.Println("Failed to migrate Version model")
		return err
	}
	err = db.AutoMigrate(&model.Stage{})
	if err != nil {
		log.Println("Failed to migrate Stage model")
		return err
	}
	err = db.AutoMigrate(&model.Environment{})
	if err != nil {
		log.Println("Failed to migrate Environment model")
		return err
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Println("Failed to migrate User model")
		return err
	}
	err = db.AutoMigrate(&model.Token{})
	if err != nil {
		log.Println("Failed to migrate Token model")
		return err
	}
	err = db.AutoMigrate(&model.UserProjectRole{})
	if err != nil {
		log.Println("Failed to migrate UserProjectRole model")
		return err
	}
	err = db.AutoMigrate(&model.Role{})
	if err != nil {
		log.Println("Failed to migrate Role model")
		return err
	}
	err = db.AutoMigrate(&model.Permission{})
	if err != nil {
		log.Println("Failed to migrate Permission model")
		return err
	}
	err = db.AutoMigrate(&model.Agent{})
	if err != nil {
		log.Println("Failed to migrate Agent model")
		return err
	}
	err = db.AutoMigrate(&model.AgentLog{})
	if err != nil {
		log.Println("Failed to migrate AgentLog model")
	}
	fmt.Printf("\nDatabase migrated\n")
	return nil
}
