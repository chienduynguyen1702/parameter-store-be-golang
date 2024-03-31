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
// @Router /api/v1/setting/roles [get]
func ListRole(c *gin.Context) {
	var roles []models.Role
	if err := DB.Preload("Permissions").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}
	type roleResponse struct {
		ID               uint                `json:"id"`
		Name             string              `json:"name"`
		Description      string              `json:"description"`
		Permissions      []models.Permission `json:"permissions"`
		PermissionsCount int                 `json:"permissions_count"`
		UserCount        int                 `json:"user_count"`
	}
	var rolesResponse []roleResponse
	for _, role := range roles {
		var userCount int64
		if err := DB.Model(&models.UserProjectRole{}).Where("role_id = ?", role.ID).Count(&userCount).Error; err != nil {
			panic("failed to count users with admin role")
		}
		rolesResponse = append(rolesResponse, roleResponse{
			ID:               role.ID,
			Name:             role.Name,
			Description:      role.Description,
			PermissionsCount: len(role.Permissions),
			Permissions:      role.Permissions,
			UserCount:        int(userCount),
		})
	}
	c.JSON(http.StatusOK, gin.H{"roles": rolesResponse})
}
