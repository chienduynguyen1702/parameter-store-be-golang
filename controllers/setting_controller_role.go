package controllers

import (
	"log"
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
// @Router /api/v1/settings/roles [get]
func ListRole(c *gin.Context) {
	org_id, exist := c.Get("org_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization id"})
		return
	}

	var roles []models.Role
	if err := DB.Preload("Permissions").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}
	type roleResponse struct {
		ID               uint   `json:"id"`
		Name             string `json:"name"`
		Description      string `json:"description"`
		PermissionsCount int    `json:"permissions_count"`
		UserCount        int    `json:"user_count"`
	}
	var rolesResponse []roleResponse
	// count user in role
	for _, role := range roles {
		// count Organization Admin
		if role.Name == "Organization Admin" {
			var orgAdminCount int64
			if err := DB.Model(&models.User{}).Where("is_organization_admin = ? AND organization_id = ? ", true, org_id).Count(&orgAdminCount).Error; err != nil {
				log.Printf("failed to count organization admins")
			}
			rolesResponse = append(rolesResponse, roleResponse{
				ID:               0,
				Name:             "Organization Admin",
				Description:      "Admin of the organization",
				PermissionsCount: len(role.Permissions),
				UserCount:        int(orgAdminCount),
			})
			continue
		}

		// count Project Admin and Developer
		var userCount int64
		if err := DB.Model(&models.UserRoleProject{}).
			Joins("left join projects on user_project_roles.project_id = projects.id").
			Where("user_project_roles.role_id = ? AND projects.organization_id = ? ", role.ID, org_id).
			Distinct("user_id").
			Count(&userCount).Error; err != nil {
			log.Printf("failed to count users with %v role", role.Name)
		}
		rolesResponse = append(rolesResponse, roleResponse{
			ID:               role.ID,
			Name:             role.Name,
			Description:      role.Description,
			PermissionsCount: len(role.Permissions),
			UserCount:        int(userCount),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"roles": rolesResponse},
	})

}
