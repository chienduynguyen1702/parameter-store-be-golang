package controllers

import (
	"net/http"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

// Get all environments
// @Summary Get all environments
// @Description Get all environments
// @Tags  Environments
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"environments": "environments"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list environments"}"
// @Router /api/v1/envs [get]
func GetEnvironments(c *gin.Context) {
	var environments []models.Environment
	if err := DB.Find(&environments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list environments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"environments": environments})
}
