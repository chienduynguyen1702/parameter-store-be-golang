package controllers

import (
	"net/http"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

// Create new stage
// @Summary Create new stage
// @Description Create new stage
// @Tags Project Detail / Parameters / Stages
// @Accept json
// @Produce json
// @Param Stage body controllers.CreateStage.createStageRequestBody true "Stage"
// @Success 201 string {string} json "{"message": "Stage created successfully", "stage": "stage"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create stage"}"
// @Router /api/v1/project/{project_id}/stages [post]
func CreateStage(c *gin.Context) {
	type createStageRequestBody struct {
		Name string `json:"name" binding:"required"`
	}
	r := createStageRequestBody{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}

	newStage := models.Stage{
		Name:      r.Name,
		ProjectID: projectID.(uint),
	}
	if err := DB.Create(&newStage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stage"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Stage created successfully", "stage": newStage})
}

// Get all stages
// @Summary Get all stages
// @Description Get all stages
// @Tags Project Detail / Parameters / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"stages": "stages"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list stages"}"
// @Router /api/v1/project/{project_id}/stages [get]
func GetStages(c *gin.Context) {
	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}

	var stages []models.Stage
	if err := DB.Find(&stages, "project_id = ?", projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list stages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stages": stages})
}

// Update stage
// @Summary Update stage
// @Description Update stage
// @Tags Project Detail / Parameters / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param stage_id path string true "Stage ID"
// @Param Stage body controllers.UpdateStage.updateStageRequestBody true "Stage"
// @Success 200 string {string} json "{"message": "Stage updated successfully", "stage": "stage"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to update stage"}"
// @Router /api/v1/project/{project_id}/stages/{stage_id} [put]
func UpdateStage(c *gin.Context) {
	type updateStageRequestBody struct {
		Name string `json:"name" binding:"required"`
	}
	r := updateStageRequestBody{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}
	stageID, exist := c.Get("stage_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stage ID from user"})
		return
	}

	var stage models.Stage
	if err := DB.First(&stage, "id = ? AND project_id = ?", stageID, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find stage"})
		return
	}

	stage.Name = r.Name
	if err := DB.Save(&stage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stage"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stage updated successfully", "stage": stage})
}

// Delete stage
// @Summary Delete stage
// @Description Delete stage
// @Tags Project Detail / Parameters / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param stage_id path string true "Stage ID"
// @Success 200 string {string} json "{"message": "Stage deleted successfully"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to delete stage"}"
// @Router /api/v1/project/{project_id}/stages/{stage_id} [delete]
func DeleteStage(c *gin.Context) {
	projectID, exist := c.Get("project_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project ID from user"})
		return
	}
	stageID, exist := c.Get("stage_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stage ID from user"})
		return
	}

	var stage models.Stage
	if err := DB.First(&stage, "id = ? AND project_id = ?", stageID, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find stage"})
		return
	}

	if err := DB.Delete(&stage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete stage"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stage deleted successfully"})
}
