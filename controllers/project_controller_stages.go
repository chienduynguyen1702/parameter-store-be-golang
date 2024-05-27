package controllers

import (
	"log"
	"net/http"
	"parameter-store-be/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Get all stages in a project
// @Summary Get all stages in a project
// @Description Get all stages in a project
// @Tags Project Detail / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"stages": "stages"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list stages"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/stages [get]
func GetListStageInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")

	// Retrieve stages from the database using the project ID
	var stages []models.Stage
	result := DB.Where("project_id = ? AND is_archived = ? ", projectID, false).Find(&stages)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"stages": stages,
		},
	})
}

// Get a stage in a project
// @Summary Get a stage in a project
// @Description Get a stage in a project
// @Tags Project Detail / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param stage_id path string true "Stage ID"
// @Success 200 string {string} json "{"stage": "stage"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get stage"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/stages/{stage_id} [get]
func GetStageInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	// Retrieve stage ID from the URL
	stageID := c.Param("stage_id")

	// Retrieve stage from the database using the project ID and stage ID
	var stage models.Stage
	result := DB.Where("project_id = ? AND id = ? AND is_archived = ? ", projectID, stageID, false).First(&stage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"stage": stage,
		},
	})
}

type stageRequestBody struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// Create a stage in a project
// @Summary Create a stage in a project
// @Description Create a stage in a project
// @Tags Project Detail / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param request body controllers.stageRequestBody true "Stage creation request"
// @Success 201 string {string} json "{"message": "Stage created successfully", "stage": "stage"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create stage"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/stages [post]
func CreateStageInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	projectIDInt, err := strconv.Atoi(projectID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	// Parse project ID to uint
	projectIDUint := uint(projectIDInt)

	// Retrieve the request body
	r := stageRequestBody{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create a new stage in the database
	newStage := models.Stage{
		Name:        r.Name,
		Description: r.Description,
		Color:       r.Color,
		ProjectID:   projectIDUint,
	}
	if err := DB.Create(&newStage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stage"})
		return
	}

	// get user from context
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// modeling user
	u := user.(models.User)
	// log stage creation
	projectLogByUser(projectIDUint, "Create Stage", "Stage created successfully", http.StatusCreated, 0, u.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Stage created successfully",
		"stage":   newStage,
	})
}

// Update a stage in a project
// @Summary Update a stage in a project
// @Description Update a stage in a project
// @Tags Project Detail / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param stage_id path string true "Stage ID"
// @Param request body controllers.stageRequestBody true "Stage update request"
// @Success 200 string {string} json "{"message": "Stage updated successfully", "stage": "stage"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to update stage"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/stages/{stage_id} [put]
func UpdateStageInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	// Retrieve stage ID from the URL
	stageID := c.Param("stage_id")

	// Retrieve the request body
	r := stageRequestBody{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the stage in the database using the project ID and stage ID
	var stage models.Stage
	result := DB.Where("project_id = ? AND id = ?", projectID, stageID).First(&stage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stage"})
		return
	}
	stage.Name = r.Name
	stage.Description = r.Description
	stage.Color = r.Color
	if err := DB.Save(&stage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stage"})
		return
	}

	// get user from context
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// modeling user
	u := user.(models.User)
	// log stage update
	projectLogByUser(stage.ProjectID, "Update Stage", "Stage updated successfully", http.StatusOK, 0, u.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Stage updated successfully",
		"stage":   stage,
	})
}

// Get all archived stages in a project
// @Summary Get all archived stages in a project
// @Description Get all archived stages in a project
// @Tags Project Detail / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 string {string} json "{"stages": "stages"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list archived stages"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/stages/archived [get]
func GetListArchivedStageInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")

	// Retrieve archived stages from the database using the project ID
	var stages []models.Stage
	result := DB.Where("project_id = ? AND is_archived = ?", projectID, true).Find(&stages)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve archived stages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"stages": stages,
		},
	})
}

// Archive a stage in a project
// @Summary Archive a stage in a project
// @Description Archive a stage in a project
// @Tags Project Detail / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param stage_id path string true "Stage ID"
// @Success 200 string {string} json "{"message": "Stage archived successfully"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to archive stage"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/stages/{stage_id}/archive [patch]
func ArchiveStageInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	// Retrieve stage ID from the URL
	stageID := c.Param("stage_id")
	// Archive the stage in the database using the project ID and stage ID
	var stage models.Stage
	result := DB.Where("project_id = ? AND id = ?", projectID, stageID).First(&stage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stage"})
		return
	}
	if stage.IsArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stage is already archived"})
		return
	}

	// get username from context
	user, exist := c.Get("user")
	if !exist {
		log.Println("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract username
	username := user.(models.User).Username

	stage.IsArchived = true
	stage.ArchivedAt = time.Now()
	stage.ArchivedBy = username
	if err := DB.Save(&stage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive stage"})
		return
	}
	// log stage archive
	projectLogByUser(stage.ProjectID, "Archive Stage", "Stage archived successfully", http.StatusOK, 0, user.(models.User).ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Stage archived successfully",
	})
}

// Unarchive a stage in a project
// @Summary Unarchive a stage in a project
// @Description Unarchive a stage in a project
// @Tags Project Detail / Stages
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param stage_id path string true "Stage ID"
// @Success 200 string {string} json "{"message": "Stage unarchived successfully"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to unarchive stage"}"
// @Security ApiKeyAuth
// @Router /api/v1/projects/{project_id}/stages/{stage_id}/unarchive [patch]
func UnarchiveStageInProject(c *gin.Context) {
	// Retrieve project ID from the URL
	projectID := c.Param("project_id")
	// Retrieve stage ID from the URL
	stageID := c.Param("stage_id")

	// Unarchive the stage in the database using the project ID and stage ID
	var stage models.Stage
	result := DB.Where("project_id = ? AND id = ?", projectID, stageID).First(&stage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stage"})
		return
	}
	if !stage.IsArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stage is already unarchived"})
		return
	}

	stage.IsArchived = false
	if err := DB.Save(&stage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unarchive stage"})
		return
	}
	// get username from context
	user, exist := c.Get("user")
	if !exist {
		log.Println("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	// Type assertion to extract username
	uID := user.(models.User).ID
	// log stage unarchive
	projectLogByUser(stage.ProjectID, "Unarchive Stage", "Stage unarchived successfully", http.StatusOK, 0, uID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Stage unarchived successfully",
	})
}
