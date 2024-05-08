package controllers

import (
	"fmt"
	"math"
	"net/http"
	"parameter-store-be/models"
	"strconv"
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
	Name              string `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	AliasName         string `gorm:"type:varchar(100)" json:"alias_name" binding:"required"`
	EstablishmentDate string `json:"establishment_date" binding:"required"`
	Description       string `gorm:"type:text" json:"description" binding:"required"`
}

// UpdateOrganizationInformation godoc
// @Summary Update organization information
// @Description Update organization information
// @Tags Organization
// @Accept json
// @Produce json
// @Param organization_id path int true "Organization ID"
// @Param Organization body organizationBody true "Organization"
// @Success 200 string {string} json "{"organizations": "organizations"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get organization"}"
// @Router /api/v1/organizations/{organization_id} [put]
func UpdateOrganizationInformation(c *gin.Context) {
	// Retrieve organization ID from the URL
	organizationIDParam := c.Param("organization_id")
	if organizationIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}
	// parse to int
	organizationIDInt, err := strconv.Atoi(organizationIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID must be an integer"})
		return
	}
	organizationID := uint(organizationIDInt)

	// Retrieve organization from the database using the user's organization ID
	var organization models.Organization
	result := DB.First(&organization, organizationID)
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
	// debug bind data
	fmt.Println(requestBody.EstablishmentDate)
	// parst string to time
	establishmentDate, err := time.Parse("02-01-2006", requestBody.EstablishmentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update organization fields
	organization.Name = requestBody.Name
	organization.AliasName = requestBody.AliasName
	organization.EstablishmentDate = establishmentDate
	organization.Description = requestBody.Description

	// Save the updated organization back to the database
	DB.Save(&organization)

	c.JSON(http.StatusOK, gin.H{"organization": organization})
}

// GetOrganizationDashboardTotals godoc
// @Summary Get organization dashboard totals
// @Description Get organization dashboard totals
// @Tags Organization
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"organizations": "organizations"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get organization"}"
// @Router /api/v1/organizations/dashboard/totals [get]
func GetOrganizationDashboardTotals(c *gin.Context) {
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
	//=======================================================
	var projectsCount int64
	DB.Model(&models.Project{}).Where("organization_id = ?", userOrganizationID).Count(&projectsCount)

	var activeProjectsCount int64
	DB.Model(&models.Project{}).Where("organization_id = ? AND is_archived != ?", userOrganizationID, true).Count(&activeProjectsCount)

	var pendingProjectsCount int64
	DB.Model(&models.Project{}).Where("organization_id = ? AND is_archived = ?", userOrganizationID, true).Count(&pendingProjectsCount)

	var usersCount int64
	DB.Model(&models.User{}).Where("organization_id = ?", userOrganizationID).Count(&usersCount)

	// count workflow in project in organization
	var totalWorkflowCount int64
	DB.Model(&models.Workflow{}).Joins("JOIN projects ON projects.id = workflows.project_id").Where("projects.organization_id = ?", userOrganizationID).Count(&totalWorkflowCount)
	//=======================================================
	var avgDurationAllWorkflowsInOrganization float64
	q := getAverageDurationByOrganizationIdQueryBuilder(userOrganizationID)
	DB.Raw(q).Row().Scan(&avgDurationAllWorkflowsInOrganization)
	roundedDuration := int(math.Round(avgDurationAllWorkflowsInOrganization))
	/*
		select count(*) from project_logs
		left join projects on project_logs.project_id = projects.id
		left join organizations on projects.organization_id = organizations.id
		where organizations.id =1
	*/
	var totalUpdatedWithinOrganization int64
	DB.Model(&models.ProjectLog{}).
		Joins("JOIN projects ON project_logs.project_id = projects.id").
		Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.id = ?", userOrganizationID).
		Count(&totalUpdatedWithinOrganization)

	var totalAgentActionsWithinOrganization int64
	DB.Model(&models.AgentLog{}).
		Joins("JOIN projects ON agent_logs.project_id = projects.id").
		Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.id = ?", userOrganizationID).
		Count(&totalAgentActionsWithinOrganization)
	firstDayOfThisMonth := get1stDayOfMonth(time.Now())
	var totalUpdatedWithinOrganizationThisMonth int64
	DB.Debug().Model(&models.ProjectLog{}).
		Joins("JOIN projects ON project_logs.project_id = projects.id").
		Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.id = ? AND project_logs.created_at BETWEEN ? AND ?", userOrganizationID, firstDayOfThisMonth, time.Now()).
		Count(&totalUpdatedWithinOrganizationThisMonth)

	var totalAgentActionsWithinOrganizationThisMonth int64
	DB.Model(&models.AgentLog{}).
		Joins("JOIN projects ON agent_logs.project_id = projects.id").
		Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.id = ? AND project_logs.created_at BETWEEN ? AND ?", userOrganizationID, firstDayOfThisMonth, time.Now()).
		Count(&totalAgentActionsWithinOrganizationThisMonth)

	// summary data
	organizationTotals := gin.H{
		"project_count":          projectsCount,
		"active_projects_count":  activeProjectsCount,
		"pending_projects_count": pendingProjectsCount,
		"user_count":             usersCount,
		"workflow_count":         totalWorkflowCount,

		"avg_duration":                   roundedDuration,
		"total_updated":                  totalUpdatedWithinOrganization,
		"total_updated_this_month":       totalUpdatedWithinOrganizationThisMonth,
		"total_agent_actions":            totalAgentActionsWithinOrganization,
		"total_agent_actions_this_month": totalAgentActionsWithinOrganizationThisMonth,
	}

	c.JSON(http.StatusOK, organizationTotals)
}

func getAverageDurationByOrganizationIdQueryBuilder(organizationID uint) string {
	return fmt.Sprintf(`
		SELECT AVG(duration) FROM workflow_logs
		JOIN projects ON projects.id = workflow_logs.project_id
		WHERE projects.organization_id = %d
	`, organizationID)
}
