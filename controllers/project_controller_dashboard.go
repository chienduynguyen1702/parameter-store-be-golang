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
	// Return the result
	c.JSON(http.StatusOK, gin.H{
		"count_updated_this_month":                countMonthParamUpdate,
		"count_agent_actions_this_month":          countMonthAgentActions,
		"count_workflows":                         len(p.Workflows),
		"count_updated_this_week":                 countWeekParamUpdate,
		"count_agent_actions_this_week":           countWeekAgentActions,
		"avg_duration_of_workflows_current_month": roundedDuration,
	})
}
func getQueryParams(c *gin.Context) (string, string, string, string) {
	granularity := c.Query("granularity")
	// fmt.Println("granularity: ", granularity) // granularity shoule be day, week, month, quarter, year
	start_date := c.Query("start_date")
	// fmt.Println("start_date: ", start_date)
	end_date := c.Query("end_date")
	// fmt.Println("end_date: ", end_date)
	workflow_id := c.Query("workflow_id")
	// fmt.Println("workflow_id: ", workflow_id)
	return granularity, start_date, end_date, workflow_id
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
// @Param project_id 	path 	string true "Project ID"
// @Param granularity 	query 	string false "Granularity: day, week, month, quarter, year, default is day"
// @Param start_date 	query 	string false "Start Date format dd-mm-yyyy"
// @Param end_date 		query 	string false "End Date format dd-mm-yyyy"
// @Param workflow_id 	query 	string false "Workflow ID specified, if not specified, get all workflows"
// @Success 200 string {string} json "{"logs": "logs"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get project dashboard logs"}"
// @Router /api/v1/projects/{project_id}/dashboard/logs [get]
func GetProjectDashboardLogs(c *gin.Context) {
	//get project_id
	projectID := c.Param("project_id")
	fmt.Println("projectID: ", projectID)

	// Get param query
	granularity, startDate, endDate, workflowID := getQueryParams(c)
	fmt.Println("granularity : ", granularity) // granularity shoule be day, week, month, quarter, year
	fmt.Println("startDate   : ", startDate)
	fmt.Println("endDate     : ", endDate)
	fmt.Println("workflowID  : ", workflowID)

	if granularity == "" {
		granularity = "day"
	}

	// type WorkflowRerunDuration struct {
	// 	WorkflowID      uint
	// 	WorkflowName    string
	// 	AverageDuration int
	// 	UnitTime        string `default:"ms"`
	// }
	type logsByGranularity struct {
		AvgDuration float64 `json:"avg_duration_in_period"`
		Count       int     `json:"count"`
		Period      string  `json:"period_start"`
		// WorkflowID  uint
	}
	// Build query
	query := queryBuilderForLogsByGranularity(granularity, startDate, workflowID)
	// fmt.Println("query: ", query)
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	// Execute query and get rows bind to struct
	var logsGranularity []logsByGranularity
	if err := DB.Raw(query).Scan(&logsGranularity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project dashboard logs"})
		return
	}
	// fmt.Println("logsGranularity: ", logsGranularity)

	// Get duration of workflow logs within current month
	// var p models.Project
	// if err := DB.Preload("Workflows").Preload("Workflows.Logs").First(&p, projectID).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
	// 	return
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
	c.JSON(http.StatusOK, gin.H{
		"logs_with_granularity": logsGranularity,
		"granularity":           granularity,
	})
}

func queryBuilderForLogsByGranularity(granularity, startDate, workflowID string) string {
	switch granularity {
	case "day": // get logs by day
		return queryBuilderForLogsByGranularityDay(startDate, workflowID)
	case "week": // get logs by week

	case "month": // get logs by month

	case "quarter": // get logs by quarter
	case "year": // get logs by year
	default:
		// get logs by day
	}
	return ""
}

func queryBuilderForLogsByGranularityDay(startDate, workflowID string) string {

	if startDate == "" {
		firstDayOfMonth := get1stDayOfMonth(time.Now()).Format("2006-01-02")
		startDate = firstDayOfMonth
	}
	return fmt.Sprintf(`
        SELECT
            to_char(date, 'YYYY-MM-DD') AS Period,
            COUNT(workflow_logs.workflow_id) AS Count,
            AVG(duration) AS Avg_Duration
        FROM
            generate_series(
            date_trunc('day', '%s'::date),
            date_trunc('day', NOW()),
            interval '1 day'
        ) AS date
        LEFT JOIN
            workflow_logs ON date_trunc('day', workflow_logs.created_at)::date = date::date
                         AND workflow_logs.state = 'completed'
                         %s
        GROUP BY
            date
        ORDER BY
            date;
    `,
		startDate,
		func() string {
			if workflowID != "" {
				return fmt.Sprintf("AND workflow_logs.workflow_id = '%s'", workflowID)
			}
			return ""
		}(),
	)
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
func get1stDayOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}
func get1stDayOfLastMonth(date time.Time) time.Time {
	return get1stDayOfMonth(date.AddDate(0, -1, 0))
}
