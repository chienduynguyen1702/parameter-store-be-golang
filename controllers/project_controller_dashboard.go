package controllers

import "github.com/gin-gonic/gin"

// GetDashboardData is a function to get dashboard data
// @Summary Get dashboard data
// @Description Get dashboard data
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"dashboard": "dashboard"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get dashboard data"}"
// @Router /api/v1/projects/{project_id}/dashboard [get]
func GetDashboardData(c *gin.Context) {
}
