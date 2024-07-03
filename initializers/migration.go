package initializers

import (
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
	err = db.AutoMigrate(&models.Version{}, &models.Project{})
	if err != nil {
		log.Println("Failed to migrate Project models")
		return err
	}
	// err = db.AutoMigrate(&models.Version{})
	// if err != nil {
	// 	log.Println("Failed to migrate Version models")
	// 	return err
	// }
	db.Migrator().CreateIndex(&models.Version{}, "agent_name_project_id")
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
	err = db.AutoMigrate(&models.AgentPullParameterLog{})
	if err != nil {
		log.Println("Failed to migrate AgentPullParameterLog models")
	}
	err = db.AutoMigrate(&models.ProjectLog{})
	if err != nil {
		log.Println("Failed to migrate ProjectLog models")
	}
	err = db.AutoMigrate(&models.Workflow{})
	if err != nil {
		log.Println("Failed to migrate Workflow models")
	}
	err = db.AutoMigrate(&models.WorkflowLog{})
	if err != nil {
		log.Println("Failed to migrate WorkflowLog models")
	}
	log.Printf("Database migrated\n")
	return nil
}
