package controllers

import (
	"net/http"
	"parameter-store-be/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProjectVersions godoc
// @Summary Get versions of project
// @Description Get versions of project
// @Tags Project Detail / Versions
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Success 200 {array} models.Version
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get versions"}"
// @Router /api/v1/projects/{project_id}/versions [get]
func GetProjectVersions(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var versions []models.Version
	DB.Where("project_id = ?", projectID).Find(&versions)
	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

// CreateNewVersion godoc
// @Summary Create new version
// @Description Create new version
// @Tags Project Detail / Versions
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Param Version body models.Version true "Version"
// @Success 200 string {string} json "{"message": "Version created"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create version"}"
// @Router /api/v1/projects/{project_id}/versions [post]
func CreateNewVersion(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var version models.Version
	if err := c.ShouldBindJSON(&version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	version.ProjectID = uint(projectID)
	DB.Create(&version)
	c.JSON(http.StatusOK, gin.H{"message": "Version created"})
}
