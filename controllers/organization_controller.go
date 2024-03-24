package controllers

import (
	"net/http"
	"parameter-store-be/model"

	"github.com/gin-gonic/gin"
)

// GetOrganization godoc
// @Summary Get organization
// @Description Get organization
// @Tags Organization
// @Accept json
// @Produce json
// @Param organization_id path string true "Organization ID"
// @Success 200 string {string} json "{"organizations": "organizations"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to get organization"}"
// @Router /api/v1/organization/{organization_id} [get]
func GetOrganization(c *gin.Context) {
	ordID := c.Param("organization_id")
	var organizations []model.Organization
	DB.First(&organizations, ordID)
	c.JSON(http.StatusOK, gin.H{"organizations": organizations})
}
