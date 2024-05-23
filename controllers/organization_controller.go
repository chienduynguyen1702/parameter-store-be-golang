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
// @Param start_date 	query 	string false "Start Date format dd-mm-yyyy"
// @Param project	 	query 	string false "Project Name specified, if not specified, get all projects"
// @Param end_date 		query 	string false "End Date format dd-mm-yyyy"
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
	// Get project name from param query
	projectName := c.Query("project")
	var cards []DashboardCard
	if projectName != "" {
		var project models.Project
		result := DB.Where("name = ? AND organization_id = ?", projectName, userOrganizationID).First(&project)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
			return
		}
		// fmt.Println("project: ", project)
		// Parse project ID to string
		projectID := project.ID
		cards = getSummaryCardsByProjectIdQueryBuilder(userOrganizationID, projectID)
	} else {
		cards = getSummaryCardsByOrganizationIdQueryBuilder(userOrganizationID)
	}
	c.JSON(http.StatusOK, cards)
}
func getSummaryCardsByOrganizationIdQueryBuilder(organizationID uint) []DashboardCard {
	// =======================================================
	var projectsCount int64
	DB.Model(&models.Project{}).Where("organization_id = ?", organizationID).Count(&projectsCount)

	var activeProjectsCount int64
	DB.Model(&models.Project{}).Where("organization_id = ? AND is_archived != ?", organizationID, true).Count(&activeProjectsCount)

	var runningAgent int64
	DB.Model(&models.Agent{}).
		Joins("LEFT JOIN projects ON projects.id = agents.project_id").
		Where("projects.organization_id = ? AND agents.is_archived = ?", organizationID, false).
		Count(&runningAgent)

	var usersCount int64
	DB.Model(&models.User{}).Where("organization_id = ?", organizationID).Count(&usersCount)

	// count workflow in project in organization
	var totalWorkflowCount int64
	DB.Model(&models.Workflow{}).Joins("JOIN projects ON projects.id = workflows.project_id").Where("projects.organization_id = ?", organizationID).Count(&totalWorkflowCount)
	// =======================================================
	var avgDurationAllWorkflowsInOrganization float64
	q := getAverageDurationByOrganizationIdQueryBuilder(organizationID)
	DB.Raw(q).Row().Scan(&avgDurationAllWorkflowsInOrganization)
	roundedDuration := int(math.Round(avgDurationAllWorkflowsInOrganization))

	var totalUpdatedWithinOrganization int64
	DB.Model(&models.ProjectLog{}).
		Joins("JOIN projects ON project_logs.project_id = projects.id").
		Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.id = ?", organizationID).
		Count(&totalUpdatedWithinOrganization)

	var totalAgentActionsWithinOrganization int64
	DB.Model(&models.AgentLog{}).
		Joins("JOIN projects ON agent_logs.project_id = projects.id").
		Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.id = ?", organizationID).
		Count(&totalAgentActionsWithinOrganization)

	var cards []DashboardCard
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-eye-red",
		Value:   int(projectsCount),
		Title:   "Projects Count",
		Name:    "project_count",
		Tooltip: "Includes all projects in organization",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-eye-red",
		Value:   int(activeProjectsCount),
		Title:   "Active Projects Count",
		Name:    "active_projects_count",
		Tooltip: "Number of running projects",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-eye-red",
		Value:   int(runningAgent),
		Title:   "Active Agent",
		Name:    "active_agent",
		Tooltip: "Number of running agents",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-eye-blue",
		Value:   int(usersCount),
		Title:   "User Count",
		Name:    "user_count",
		Tooltip: "Total active users",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-eye-red",
		Value:   int(totalWorkflowCount),
		Title:   "Workflow Count",
		Name:    "workflow_count",
		Tooltip: "Total active workflows in all projects",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-info-green",
		Value:   roundedDuration,
		Title:   "Average Duration",
		Name:    "avg_duration",
		Tooltip: "Average duration all CICD workflows in organization",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-shopping-carts-purple",
		Value:   int(totalUpdatedWithinOrganization),
		Title:   "Total Updated",
		Name:    "total_updated",
		Tooltip: "Number of updates",
	})
	//	cards = append(cards, DashboardCard{
	//		Icon:    "dashboard-shopping-carts-purple",
	//		Value:   int(totalUpdatedWithinOrganizationThisMonth),
	//		Title:   "Total Updated This Month",
	//		Name:    "total_updated_this_month",
	//		Tooltip: "Number of updates this month",
	//	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-shopping-carts-purple",
		Value:   int(totalAgentActionsWithinOrganization),
		Title:   "Total Agent Actions",
		Name:    "total_agent_actions",
		Tooltip: "Number of agent pull parameter",
	})
	return cards
}

func getSummaryCardsByProjectIdQueryBuilder(organizationID, projectID uint) []DashboardCard {
	// =======================================================
	var runningAgent int64
	DB.Model(&models.Agent{}).
		Joins("LEFT JOIN projects ON projects.id = agents.project_id").
		Where("projects.organization_id = ? AND projects.id = ? AND agents.is_archived = ?", organizationID, projectID, false).
		Count(&runningAgent)

	var usersCount int64
	DB.Model(&models.UserRoleProject{}).Where("project_id = ?", projectID).Count(&usersCount)

	// count workflow in project in organization
	var totalWorkflowCount int64
	DB.Model(&models.Workflow{}).Joins("JOIN projects ON projects.id = workflows.project_id").Where("projects.id = ?", projectID).Count(&totalWorkflowCount)
	// =======================================================
	var avgDurationAllWorkflowsInOrganization float64
	// q := getAverageDurationByOrganizationIdQueryBuilder(organizationID)
	DB.
		Select("AVG(duration)").
		Table("workflow_logs").
		Joins("JOIN projects ON projects.id = workflow_logs.project_id").
		Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.id = ? AND projects.id = ? ", organizationID, projectID).
		Row().Scan(&avgDurationAllWorkflowsInOrganization)
	roundedDuration := int(math.Round(avgDurationAllWorkflowsInOrganization))

	var totalUpdatedWithinOrganization int64
	DB.Model(&models.ProjectLog{}).
		Joins("JOIN projects ON project_logs.project_id = projects.id").
		Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.id = ? AND projects.id = ? ", organizationID, projectID).
		Count(&totalUpdatedWithinOrganization)

	var totalAgentActionsWithinOrganization int64
	DB.Model(&models.AgentLog{}).
		Joins("JOIN projects ON agent_logs.project_id = projects.id").
		Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.id = ? AND projects.id = ? ", organizationID, projectID).
		Count(&totalAgentActionsWithinOrganization)

	var cards []DashboardCard
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-eye-red",
		Value:   int(runningAgent),
		Title:   "Active Agent",
		Name:    "active_agent",
		Tooltip: "Number of running agents",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-eye-blue",
		Value:   int(usersCount),
		Title:   "User Count",
		Name:    "user_count",
		Tooltip: "Total active users",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-eye-red",
		Value:   int(totalWorkflowCount),
		Title:   "Workflow Count",
		Name:    "workflow_count",
		Tooltip: "Total active workflows in all projects",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-info-green",
		Value:   roundedDuration,
		Title:   "Average Duration",
		Name:    "avg_duration",
		Tooltip: "Average duration all CICD workflows in organization",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-shopping-carts-purple",
		Value:   int(totalUpdatedWithinOrganization),
		Title:   "Total Updated",
		Name:    "total_updated",
		Tooltip: "Number of updates",
	})
	cards = append(cards, DashboardCard{
		Icon:    "dashboard-shopping-carts-purple",
		Value:   int(totalAgentActionsWithinOrganization),
		Title:   "Total Agent Actions",
		Name:    "total_agent_actions",
		Tooltip: "Number of agent pull parameter",
	})
	return cards
}

func getAverageDurationByOrganizationIdQueryBuilder(organizationID uint) string {
	return fmt.Sprintf(`
		SELECT AVG(duration) FROM workflow_logs
		JOIN projects ON projects.id = workflow_logs.project_id
		WHERE projects.organization_id = %d
	`, organizationID)
}

// GetOrganizationDashboardLogs godoc
// @Summary Get organization dashboard logs
// @Description Get organization dashboard logs
// @Tags Organization
// @Accept json
// @Produce json
// @Param granularity 	query 	string false "Granularity: day, week, month, quarter, year, default is day"
// @Param start_date 	query 	string false "Start Date format dd-mm-yyyy"
// @Param project	 	query 	string false "Project Name specified, if not specified, get all projects"
// @Param end_date 		query 	string false "End Date format dd-mm-yyyy"
// @Success 200 string {string} json "{"organizations": "organizations"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get organization"}"
// @Router /api/v1/organizations/dashboard/logs [get]
func GetOrganizationDashboardLogs(c *gin.Context) {
	// get user from context
	userinContext, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	// get organization id from user in string
	userOrganizationID := strconv.Itoa(int(userinContext.OrganizationID))

	// Get param query
	granularity, from, to, _ := getQueryParams(c)
	// fmt.Println("granularity : ", granularity) // granularity shoule be day, week, month, quarter, year
	// fmt.Println("from   : ", from)
	// fmt.Println("to     : ", to)
	// fmt.Println("workflowID  : ", workflowID)
	if granularity == "" {
		granularity = "day"
	}
	type logsByGranularity struct {
		AvgDuration float64 `json:"avg_duration_in_period"`
		Count       int     `json:"count"`
		Period      string  `json:"period_start"`
		// WorkflowID  uint
	}
	// Get project name from param query
	projectName := c.Query("project")
	var query string
	if projectName != "" {
		var project models.Project
		result := DB.Where("name = ? AND organization_id = ?", projectName, userOrganizationID).First(&project)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
			return
		}
		// fmt.Println("project: ", project)
		// Parse project ID to string
		projectID := strconv.Itoa(int(project.ID))
		query = queryBuilderForLogsByGranularityOfProject(granularity, from, to, projectID)
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
	} else {
		// Build query
		query = queryBuilderForLogsByGranularityOfOrganization(granularity, from, to, userOrganizationID)
		// fmt.Println("query: ", query)
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
	}
	// Execute query and get rows bind to struct
	var logsGranularity []logsByGranularity
	if err := DB.Raw(query).Scan(&logsGranularity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project dashboard logs"})
		return
	}
	// fmt.Println("logsGranularity: ", logsGranularity)

	c.JSON(http.StatusOK, gin.H{
		"logs_with_granularity": logsGranularity,
		"granularity":           granularity,
	})

}

func queryBuilderForLogsByGranularityOfOrganization(granularity, from, to, organizationID string) string {
	if from == "" {
		switch granularity {
		case "day":
			from = time.Now().AddDate(0, 0, -14).Format("2006-01-02")
		case "week":
			from = getDateTimeOfMondayOfWeek(time.Now().AddDate(0, 0, -42)).Format("2006-01-02")
		case "month", "quarter", "year":
			from = get1stDayOfYear(time.Now()).Format("2006-01-02")
		}
	}
	if to == "" {
		to = time.Now().Format("2006-01-02")
	}

	dateTrunc, interval := "", ""
	switch granularity {
	case "day":
		dateTrunc, interval = "day", "1 day"
	case "week":
		dateTrunc, interval = "week", "1 week"
	case "month":
		dateTrunc, interval = "month", "1 month"
	case "quarter":
		dateTrunc, interval = "quarter", "3 month"
	case "year":
		dateTrunc, interval = "year", "1 year"
	}

	return fmt.Sprintf(`
		SELECT
			to_char(date, 'YYYY-MM-DD') AS Period,
			COUNT(workflow_logs.workflow_id) AS Count,
			AVG(duration) AS Avg_Duration
		FROM
			generate_series(
			date_trunc('%s', '%s'::date),
			date_trunc('%s', '%s'::date),
			interval '%s'
		) AS date
		LEFT JOIN
			workflow_logs ON date_trunc('%s', workflow_logs.created_at)::date = date::date
			AND workflow_logs.state = 'completed'
			AND workflow_logs.project_id IN (
				SELECT
					id
				FROM
					projects
				WHERE
					projects.organization_id = %s
			)
		GROUP BY
			date
		ORDER BY
			date;
	`,
		dateTrunc, from,
		dateTrunc, to,
		interval,
		dateTrunc,
		organizationID,
	)
}

func queryBuilderForLogsByGranularityOfProject(granularity, from, to, projectID string) string {
	if from == "" {
		switch granularity {
		case "day":
			from = time.Now().AddDate(0, 0, -14).Format("2006-01-02")
		case "week":
			from = getDateTimeOfMondayOfWeek(time.Now().AddDate(0, 0, -42)).Format("2006-01-02")
		case "month", "quarter", "year":
			from = get1stDayOfYear(time.Now()).Format("2006-01-02")
		}
	}
	if to == "" {
		to = time.Now().Format("2006-01-02")
	}

	dateTrunc, interval := "", ""
	switch granularity {
	case "day":
		dateTrunc, interval = "day", "1 day"
	case "week":
		dateTrunc, interval = "week", "1 week"
	case "month":
		dateTrunc, interval = "month", "1 month"
	case "quarter":
		dateTrunc, interval = "quarter", "3 month"
	case "year":
		dateTrunc, interval = "year", "1 year"
	}

	return fmt.Sprintf(`
		SELECT
			to_char(date, 'YYYY-MM-DD') AS Period,
			COUNT(workflow_logs.workflow_id) AS Count,
			AVG(duration) AS Avg_Duration
		FROM
			generate_series(
			date_trunc('%s', '%s'::date),
			date_trunc('%s', '%s'::date),
			interval '%s'
		) AS date
		LEFT JOIN
			workflow_logs ON date_trunc('%s', workflow_logs.created_at)::date = date::date
			AND workflow_logs.state = 'completed'
			AND workflow_logs.project_id = %s
		GROUP BY
			date
		ORDER BY
			date;
	`,
		dateTrunc, from,
		dateTrunc, to,
		interval,
		dateTrunc,
		projectID,
	)
}
