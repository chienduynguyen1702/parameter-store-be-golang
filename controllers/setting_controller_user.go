package controllers

import (
	"net/http"
	"parameter-store-be/models"

	"github.com/gin-gonic/gin"
)

// ListUser godoc
// @Summary List users
// @Description List users
// @Tags Setting / User
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"users": "users"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list users"}"
// @Router /api/v1/settings/users [get]
func ListUser(c *gin.Context) {
	claims, err := parseJWTTokenFromCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	org_id := claims["org_id"]
	var users []models.User
	DB.Where("organization_id = ? AND is_archived != ?", org_id, true).Find(&users)

	type userResponse struct {
		ID                  uint   `json:"id"`
		Email               string `json:"email"`
		Username            string `json:"username"`
		Phone               string `json:"phone"`
		IsOrganizationAdmin bool   `json:"is_organization_admin"`
	}
	var usersResponse []userResponse
	for _, user := range users {
		usersResponse = append(usersResponse, userResponse{
			ID:                  user.ID,
			Email:               user.Email,
			Username:            user.Username,
			Phone:               user.Phone,
			IsOrganizationAdmin: user.IsOrganizationAdmin,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"users": usersResponse,
	}})
}

// CreateUser godoc
// @Summary Create user
// @Description Create user
// @Tags Setting / User
// @Accept json
// @Produce json
// @Param User body controllers.CreateUser.createUserRequestBody true "User"
// @Success 201 string {string} json "{"message": "User created successfully", "user": "user"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to create user"}"
// @Router /api/v1/settings/users [post]
func CreateUser(c *gin.Context) {
	type createUserRequestBody struct {
		Email    string `json:"email" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Phone    string `json:"phone"`
	}
	r := createUserRequestBody{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	org_id, exist := c.Get("org_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization ID from user"})
		return
	}
	hash, err := generateBcryptPassword(r.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUser := models.User{
		Email:               r.Email,
		Username:            r.Username,
		OrganizationID:      org_id.(uint),
		Password:            string(hash),
		Phone:               r.Phone,
		IsOrganizationAdmin: false,
	}
	if err := DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": newUser})
}

// UpdateUserInformation godoc
// @Summary Update user information
// @Description Update user information
// @Tags Setting / User
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param User body controllers.UpdateUserInformation.updateUserRequestBody true "User"
// @Success 200 string {string} json "{"message": "User information updated successfully", "user": "user"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to update user information"}"
// @Router /api/v1/settings/users/{user_id} [put]
func UpdateUserInformation(c *gin.Context) {
	type updateUserRequestBody struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Phone    string `json:"phone"`
		Passowrd string `json:"password"`
	}
	r := updateUserRequestBody{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_id := c.Param("user_id")
	var user models.User
	if err := DB.First(&user, user_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	hashPassword, err := generateBcryptPassword(r.Passowrd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Email = r.Email
	user.Username = r.Username
	user.Phone = r.Phone
	user.Password = hashPassword
	if err := DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user information"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User information updated successfully", "user": user})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user including all its data in user_project_role table
// @Tags Setting / User
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 string {string} json "{"message": "User deleted"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to delete user"}"
// @Router /api/v1/settings/users/{user_id} [delete]
func DeleteUser(c *gin.Context) {
	user_id := c.Param("user_id")
	var user models.User
	// check if user exists
	if err := DB.First(&user, user_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Delete user from user_project_role table
	if err := DB.Where("user_id = ?", user_id).Delete(&models.UserProjectRole{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user from user_project_role table"})
		return
	}

	// Delete user from user table
	if err := DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// ListArchivedUser godoc
// @Summary List archived users
// @Description List archived users
// @Tags Setting / User
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"users": "users"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to list archived users"}"
// @Router /api/v1/settings/users/archived [get]
func ListArchivedUser(c *gin.Context) {
	claims, err := parseJWTTokenFromCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	org_id := claims["org_id"]
	var users []models.User
	DB.Where("organization_id = ? AND is_archived = ?", org_id, true).Find(&users)

	type userResponse struct {
		ID         uint   `json:"id"`
		Email      string `json:"email"`
		Username   string `json:"username"`
		Phone      string `json:"phone"`
		AvatarURL  string `json:"avatar_url"`
		ArchivedAt string `json:"archived_at"`
		ArchivedBy string `json:"archived_by"`
		// IsOrganizationAdmin bool   `json:"is_organization_admin"`
	}
	var usersResponse []userResponse
	for _, user := range users {
		usersResponse = append(usersResponse, userResponse{
			ID:         user.ID,
			Email:      user.Email,
			Username:   user.Username,
			AvatarURL:  user.AvatarURL,
			Phone:      user.Phone,
			ArchivedBy: user.ArchivedBy,
			ArchivedAt: user.ArchivedAt.String(),
			// IsOrganizationAdmin: user.IsOrganizationAdmin,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"users": usersResponse,
		}})
}

// ArchiveUser godoc
// @Summary Archive user
// @Description Archive user
// @Tags Setting / User
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 string {string} json "{"message": "User archived"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to archive user"}"
// @Router /api/v1/settings/users/{user_id}/archive [put]
func ArchiveUser(c *gin.Context) {
	user_id := c.Param("user_id")
	var user models.User
	// check if user exists
	if err := DB.First(&user, user_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	user.IsArchived = true
	if err := DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User archived"})
}

// RestoreUser godoc
// @Summary  Restore archived user
// @Description Restore archived user
// @Tags Setting / User
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 string {string} json "{"message": "User restored"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to unarchive user"}"
// @Router /api/v1/settings/users/{user_id}/restore [put]
func RestoreUser(c *gin.Context) {
	user_id := c.Param("user_id")
	var user models.User
	// check if user exists
	if err := DB.First(&user, user_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	user.IsArchived = false
	if err := DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unarchive user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User restored"})
}
