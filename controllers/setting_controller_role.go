package controllers

import (
	"net/http"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

// ListRole godoc
// @Summary List roles
// @Description List roles and its permissions
// @Tags Setting / Role
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"roles": "roles"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list roles"}"
// @Router /api/v1/setting/role [get]
func ListRole(c *gin.Context) {
	var roles []models.Role
	if err := DB.Preload("Permissions").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}
