package controllers

import (
	"fmt"
	"net/http"
	"parameter-store-be/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GetOrganizationInformation godoc
// @Summary Get organization information
// @Description Get organization information
// @Tags Organization
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"organizations": "organizations"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get organization"}"
// @Router /api/v1/organizations/ [get]
func GetOrganizationInformation(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	userOrganizationID := user.(models.User).OrganizationID
	if userOrganizationID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization ID from user"})
		return
	}

	var organization models.Organization
	result := DB.First(&organization, userOrganizationID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organization"})
		return
	}

	var usersCount int64
	var projectsCount int64

	DB.Model(&models.User{}).Where("organization_id = ?", userOrganizationID).Count(&usersCount)
	DB.Model(&models.Project{}).Where("organization_id = ?", userOrganizationID).Count(&projectsCount)

	organizationResponseBody := gin.H{
		"organization": gin.H{
			"id":                 organization.ID,
			"create_at":          organization.CreatedAt,
			"name":               organization.Name,
			"alias_name":         organization.AliasName,
			"establishment_date": organization.EstablishmentDate,
			"description":        organization.Description,
			"user_count":         usersCount,
			"project_count":      projectsCount,
			"address":            organization.Address,
		},
	}

	c.JSON(http.StatusOK, organizationResponseBody)
}

type organizationBody struct {
	Name              string `gorm:"type:varchar(100);not null"`
	AliasName         string `gorm:"type:varchar(100)"`
	EstablishmentDate time.Time
	Description       string `gorm:"type:text"`
}

// UpdateOrganizationInformation godoc
// @Summary Update organization information
// @Description Update organization information
// @Tags Organization
// @Accept json
// @Produce json
// @Param Organization body organizationBody true "Organization"
// @Success 200 string {string} json "{"organizations": "organizations"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get organization"}"
// @Router /api/v1/organizations/ [put]
func UpdateOrganizationInformation(c *gin.Context) {
	// Retrieve user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract organization ID
	userOrganizationID := user.(models.User).OrganizationID
	if userOrganizationID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization ID from user"})
		return
	}
	fmt.Println("debug")
	// Retrieve organization from the database using the user's organization ID
	var organization models.Organization
	result := DB.First(&organization, userOrganizationID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organization"})
		return
	}

	// Bind JSON data to organizationBody struct
	var requestBody organizationBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update organization fields
	organization.Name = requestBody.Name
	organization.AliasName = requestBody.AliasName
	organization.EstablishmentDate = requestBody.EstablishmentDate
	organization.Description = requestBody.Description

	// Save the updated organization back to the database
	DB.Save(&organization)

	c.JSON(http.StatusOK, gin.H{"organization": organization})
}
