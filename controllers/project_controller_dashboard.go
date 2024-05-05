package controllers

import (
	"fmt"
	"math"
	"net/http"
	"parameter-store-be/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GetDashboardData is a function to get dashboard data
// @Summary Get dashboard data
// @Description Get dashboard data
// @Tags Project Detail / Dashboard
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"dashboard": "dashboard"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get dashboard data"}"
// @Router /api/v1/projects/{project_id}/dashboard/totals [get]
func GetProjectDashboardTotals(c *gin.Context) {
	//get project_id
	projectID := c.Param("project_id")
	// count update within current month in project_logs

	// Get current month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)
	// count updated logs within current month
	countMonthParamUpdate := countMonthParamUpdate(c, projectID, startOfMonth, endOfMonth)
	if countMonthParamUpdate < 0 {
		return
	}
	// count agent actions within current month
	countMonthAgentActions := countMonthAgentActions(c, projectID, startOfMonth, endOfMonth)
	if countMonthAgentActions < 0 {
		return
	}

	// Get current week
	startOfWeek := getDateTimeOfMondayOfWeek(now)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	// log.Println(startOfWeek, endOfWeek)

	// count updated logs within current week
	countWeekParamUpdate := countWeekParamUpdate(c, projectID, startOfWeek, endOfWeek)
	if countWeekParamUpdate < 0 {
		return
	}
	// Get count of agent logs within current week
	countWeekAgentActions := countWeekAgentActions(c, projectID, startOfWeek, endOfWeek)
	if countWeekAgentActions < 0 {
		return
	}

	// Get duration of workflow logs within current month
	var p models.Project
	if err := DB.Preload("Workflows").Preload("Workflows.Logs").First(&p, projectID).Error; err != nil {
		// if err := DB.Preload("Workflows", "started_at BETWEEN ? AND ?", startOfMonth, endOfMonth).First(&p, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}

	// calculate avg_duration_of_workflows_current_month
	var avgDurationAllWorkflows float64
	var listWorkflowIds []uint
	for _, workflow := range p.Workflows {
		listWorkflowIds = append(listWorkflowIds, workflow.WorkflowID)
	}
	DB.Model(&models.WorkflowLog{}).Where("workflow_id IN (?)", listWorkflowIds).Select("AVG(duration)").Row().Scan(&avgDurationAllWorkflows)

	roundedDuration := int(math.Round(avgDurationAllWorkflows))
	// calculate avg_duration_of_workflows_current_month
	// type WorkflowRerunDuration struct {
	// 	WorkflowID      uint
	// 	WorkflowName    string
	// 	AverageDuration int
	// 	UnitTime        string `default:"ms"`
	// }
	// var wrds []WorkflowRerunDuration
	// var workflowLogInThisProject []models.WorkflowLog
	// // get average duration of workflow logs
	// for _, workflow := range p.Workflows {

	// 	// log.Println("Workflow: ", workflow)
	// 	var avgDuration float64
	// 	DB.Model(&models.WorkflowLog{}).Where("workflow_id = ?", workflow.WorkflowID).Select("AVG(duration)").Row().Scan(&avgDuration)

	// 	roundedDuration := int(math.Round(avgDuration))
	// 	// fmt.Printf("WorkflowID: %d, Average Duration: %d\n", workflow.WorkflowID, roundedDuration)

	// 	wrds = append(wrds, WorkflowRerunDuration{
	// 		WorkflowID:      workflow.WorkflowID,
	// 		WorkflowName:    workflow.Name,
	// 		AverageDuration: roundedDuration,
	// 		UnitTime:        "ms",
	// 	})
	// 	workflowLogInThisProject = append(workflowLogInThisProject, workflow.Logs...)
	// }
	// log.Println("wrds: ", wrds)
	// Return the result
	c.JSON(http.StatusOK, gin.H{
		"count_updated_this_month":                countMonthParamUpdate,
		"count_agent_actions_this_month":          countMonthAgentActions,
		"count_workflows":                         len(p.Workflows),
		"count_updated_this_week":                 countWeekParamUpdate,
		"count_agent_actions_this_week":           countWeekAgentActions,
		"avg_duration_of_workflows_current_month": roundedDuration,
		// "logs":                           workflowLogInThisProject,
	})
}

func countMonthParamUpdate(c *gin.Context, projectID string, startOfMonth time.Time, endOfMonth time.Time) int64 {
	var countMonthUpdate int64
	if err := DB.Model(&models.ProjectLog{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startOfMonth, endOfMonth).
		Count(&countMonthUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count project logs"})
		return -1
	}
	return countMonthUpdate
}
func countMonthAgentActions(c *gin.Context, projectID string, startOfMonth time.Time, endOfMonth time.Time) int64 {
	var countMonthAgentActions int64
	if err := DB.Model(&models.AgentLog{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startOfMonth, endOfMonth).
		Count(&countMonthAgentActions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count agent logs"})
		return -1
	}
	return countMonthAgentActions
}
func countWeekParamUpdate(c *gin.Context, projectID string, startOfWeek time.Time, endOfWeek time.Time) int64 {
	var countWeekUpdate int64
	if err := DB.Model(&models.ProjectLog{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startOfWeek, endOfWeek).
		Count(&countWeekUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count project logs"})
		return -1
	}
	return countWeekUpdate
}
func countWeekAgentActions(c *gin.Context, projectID string, startOfWeek time.Time, endOfWeek time.Time) int64 {
	var countWeekAgentActions int64
	if err := DB.Model(&models.AgentLog{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startOfWeek, endOfWeek).
		Count(&countWeekAgentActions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count agent logs"})
		return -1
	}
	return countWeekAgentActions
}

// GetProjectDashboardLogs godoc
// @Summary Get project dashboard logs
// @Description Get project dashboard logs
// @Tags Project Detail / Dashboard
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"logs": "logs"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get project dashboard logs"}"
// @Router /api/v1/projects/{project_id}/dashboard/logs [get]
func GetProjectDashboardLogs(c *gin.Context) {
	//get project_id
	projectID := c.Param("project_id")

	// Get granularity from query
	granularity := c.Query("granularity")
	fmt.Println("granularity: ", granularity) // granularity shoule be day, week, month, quarter, year

	// Get duration of workflow logs within current month
	var p models.Project
	if err := DB.Preload("Workflows").Preload("Workflows.Logs").First(&p, projectID).Error; err != nil {
		// if err := DB.Preload("Workflows", "started_at BETWEEN ? AND ?", startOfMonth, endOfMonth).First(&p, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	type WorkflowRerunDuration struct {
		WorkflowID      uint
		WorkflowName    string
		AverageDuration int
		UnitTime        string `default:"ms"`
	}
	var wrds []WorkflowRerunDuration
	var workflowLogInThisProject []models.WorkflowLog
	// get average duration of workflow logs
	for _, workflow := range p.Workflows {

		// log.Println("Workflow: ", workflow)
		var avgDuration float64
		DB.Model(&models.WorkflowLog{}).Where("workflow_id = ?", workflow.WorkflowID).Select("AVG(duration)").Row().Scan(&avgDuration)

		roundedDuration := int(math.Round(avgDuration))
		// fmt.Printf("WorkflowID: %d, Average Duration: %d\n", workflow.WorkflowID, roundedDuration)

		wrds = append(wrds, WorkflowRerunDuration{
			WorkflowID:      workflow.WorkflowID,
			WorkflowName:    workflow.Name,
			AverageDuration: roundedDuration,
			UnitTime:        "ms",
		})
		workflowLogInThisProject = append(workflowLogInThisProject, workflow.Logs...)
	}
	c.JSON(http.StatusOK, gin.H{"logs": workflowLogInThisProject})
}

func queryBuilderForLogsByGranularity(granularity string) string {
	switch granularity {
	case "day": // get logs by day

	case "week": // get logs by week

	case "month": // get logs by month

	case "quarter": // get logs by quarter
	case "year": // get logs by year
	default:
		// get logs by day
	}
	return ""
}

// getLogsByGranularity is a function to get logs by granularity
func getLogsByGranularity(granularity string) string {
	switch granularity {
	case "day":

	case "week": // get logs by week

	case "month": // get logs by month

	case "quarter": // get logs by quarter
	case "year": // get logs by year
	default:
		// get logs by day
	}
	return ""
}

// get Date Of Monday Of Week for a given date
func getDateTimeOfMondayOfWeek(date time.Time) time.Time {

	// get the day of the week for the given date
	dayOfWeek := date.Weekday()

	// get the number of days to subtract from the given date to get the first day of the week
	daysToSubtract := int(dayOfWeek) - 1

	// if the day of the week is Sunday, then subtract 6 days to get the first day of the week
	if dayOfWeek == time.Sunday {
		daysToSubtract = 6
	}

	// get the first day of the week
	firstDayOfWeek := date.AddDate(0, 0, -daysToSubtract)
	firstDayOfWeek = time.Date(firstDayOfWeek.Year(), firstDayOfWeek.Month(), firstDayOfWeek.Day(), 0, 0, 0, 0, firstDayOfWeek.Location())

	return firstDayOfWeek
}
