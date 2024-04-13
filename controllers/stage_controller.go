package controllers

import (
	"net/http"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

// Get all stages
// @Summary Get all stages
// @Description Get all stages
// @Tags  Stages
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"stages": "stages"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list stages"}"
// @Router /api/v1/stages [get]
func GetStages(c *gin.Context) {
	var stages []models.Stage
	if err := DB.Find(&stages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list stages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stages": stages})
}
