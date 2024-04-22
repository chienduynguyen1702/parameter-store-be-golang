package controllers

import (
	"log"
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
// @Router /api/v1/projects/{project_id}/dashboard [get]
func GetProjectDashboard(c *gin.Context) {
	//get project_id
	projectID := c.Param("project_id")
	// count update within current month in project_logs

	// Get current month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var countMonthUpdate int64
	if err := DB.Model(&models.ProjectLog{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startOfMonth, endOfMonth).
		Count(&countMonthUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count project logs"})
		return
	}
	var countMonthAgentActions int64
	if err := DB.Model(&models.AgentLog{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startOfMonth, endOfMonth).
		Count(&countMonthAgentActions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count agent logs"})
		return
	}
	// Get current week
	startOfWeek := getDateTimeOfMondayOfWeek(now)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	log.Println(startOfWeek, endOfWeek)
	// Get count of project logs within current week
	var countWeekUpdate int64
	if err := DB.Model(&models.ProjectLog{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startOfWeek, endOfWeek).
		Count(&countWeekUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count project logs"})
		return
	}
	// Get count of agent logs within current week
	var countWeekAgentActions int64
	if err := DB.Model(&models.AgentLog{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startOfWeek, endOfWeek).
		Count(&countWeekAgentActions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count agent logs"})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{
		"update_count_current_month":  countMonthUpdate,
		"agent_actions_current_month": countMonthAgentActions,
		"update_count_current_week":   countWeekUpdate,
		"agent_actions_current_week":  countWeekAgentActions,
	})
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
