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
// @Security ApiKeyAuth
// @Failure 500 string {string} json "{"error": "Failed to get versions"}"
// @Router /api/v1/projects/{project_id}/versions [get]
func GetProjectVersions(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var versions []models.Version // order by number desc
	DB.Order("number desc").Preload("Parameters").Where("project_id = ?", projectID).Find(&versions)
	// DB.Preload("Parameters").Where("project_id = ?", projectID).Find(&versions)
	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

// CreateNewVersion godoc
// @Summary Create new version
// @Description Create new version
// @Tags Project Detail / Versions
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Param versionName body controllers.CreateNewVersion.versionName true "Version name"
// @Success 200 string {string} json "{"message": "Version created"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Security ApiKeyAuth
// @Failure 500 string {string} json "{"error": "Failed to create version"}"
// @Router /api/v1/projects/{project_id}/versions [post]
func CreateNewVersion(c *gin.Context) {
	type versionName struct {
		Number      string `json:"release_version"`
		Description string `json:"description"`
	}
	var v versionName
	projectID, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	uintProjectID := uint(projectID)

	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check unique name
	var count int64
	DB.Model(&models.Version{}).Where("project_id = ? AND number = ?", projectID, v.Number).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Version name already exists"})
		return
	}
	// get lastest version of project model
	var project models.Project
	if err := DB.Preload("LatestVersion").Preload("LatestVersion.Parameters").First(&project, uintProjectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}
	// copy latest version to new version, except name
	newVersion := models.Version{
		// ProjectID:   uintProjectID,
		Number:      v.Number,
		Description: v.Description,
		Name:        v.Number,
	}

	// Create association records for parameters
	// newVersion.Parameters = append(newVersion.Parameters, project.LatestVersion.Parameters...)

	if err := DB.Create(&newVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create version"})
		return
	}
	// Clone parameters to new version
	for _, param := range project.LatestVersion.Parameters {
		if param.IsArchived {
			continue
		}
		newParam := models.Parameter{
			StageID:       param.StageID,
			EnvironmentID: param.EnvironmentID,
			Name:          param.Name,
			Value:         param.Value,
			Description:   param.Description,
			ProjectID:     param.ProjectID,
			IsArchived:    param.IsArchived,
			ArchivedBy:    param.ArchivedBy,
			ArchivedAt:    param.ArchivedAt,
			IsApplied:     param.IsApplied,
			EditedAt:      param.EditedAt,
		}
		newVersion.Parameters = append(newVersion.Parameters, newParam)
	}
	if err := DB.Save(&newVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clone parameter to new version"})
		return
	}
	// update project latest version
	project.LatestVersionID = newVersion.ID
	if err := DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project latest version"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Version created"})
}
