package controllers

import (
	"net/http"
	"os"
	"parameter-store-be/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user in the database

// Register godoc
// @Summary Register a new user and organization
// @Description Register a new user and organization
// @Tags Auth
// @Accept json
// @Produce json
// @Param Creadentials body controllers.Register.registerRequestBody true "User registration request"
// @Success 201 string {string} json "{"message": "User registered successfully", "user": {email: "	email", organization_id: "organization_id"}}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to register user"}"
// @Router /api/v1/auth/register [post]
func Register(c *gin.Context) {
	type registerRequestBody struct {
		Email            string `json:"email" binding:"required"`
		Password         string `json:"password" binding:"required"`
		OrganizationName string `json:"organization_name" binding:"required"`
	}
	r := registerRequestBody{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := generateBcryptPassword(r.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newOrganization := models.Organization{
		Name: r.OrganizationName,
	}
	if err := DB.Create(&newOrganization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register organization"})
		return
	}
	organizationID := newOrganization.ID

	newUser := models.User{
		Email:               r.Email,
		Username:            r.Email,
		Password:            string(hash),
		OrganizationID:      organizationID,
		IsOrganizationAdmin: true,
	}
	// Create the user in the database
	if err := DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	// Create the organization in the database

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "email:": r.Email, "organization_name:": r.OrganizationName})

}

// Login logs in a user, if successful, set cookie header
// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body controllers.Login.loginRequestBody true "User login request"
// @Success 200 string {string} json "{"message": "User logged in successfully", "user": {email: "email", organization_id: "organization_id"}}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 401 string {string} json "{"error": "Unauthorized"}"
// @Failure 500 string {string} json "{"error": "Failed to login user"}"
// @Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	type loginRequestBody struct {
		Email            string `json:"email" binding:"required"`
		Password         string `json:"password" binding:"required"`
		OrganizationName string `json:"organization_name" binding:"required"`
	}
	l := loginRequestBody{}
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	organization := models.Organization{
		Name: l.OrganizationName,
	}

	// Check if the organization exists in the database
	if err := DB.Where("name = ?", l.OrganizationName).First(&organization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	var user models.User
	// Check if the user exists in the database by email and organization_name
	if err := DB.Where("email = ? AND organization_id = ?", l.Email, organization.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	// Check if the user is archived
	if user.IsArchived {
		c.SetCookie("Authorization", "", -1, "/", "", false, true)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}
	var responseLogedInUser struct {
		Username       string `json:"username"`
		Email          string `json:"email"`
		OrganizationID uint   `json:"organization_id"`
	}
	responseLogedInUser.Username = user.Username
	responseLogedInUser.Email = user.Email
	responseLogedInUser.OrganizationID = user.OrganizationID
	// Generate a JWT token
	jwtToken, err := generateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token for user"})
		return
	}
	// set login time
	user.LastLogin = time.Now()
	DB.Save(&user)
	// Set the JWT token in a cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		"Authorization",
		jwtToken,
		3600*24*30,
		"",
		os.Getenv("COOKIE_DOMAIN"),
		true,
		true,
	)
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "status:": "success", "user": responseLogedInUser})

}

// Validate validates a user by cookie
// Validate godoc
// @Summary Validate a user
// @Description Validate a user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"message": "User logged in successfully"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 401 string {string} json "{"error": "Failed to validate user"}"
// @Failure 500 string {string} json "{"error": "Internal server error"}"
// @Router /api/v1/auth/validate [get]
func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to validate user"})
		return
	}
	// set login time
	validatedUser := user.(models.User)

	var responseLogedInUser struct {
		Username       string `json:"username"`
		Email          string `json:"email"`
		OrganizationID uint   `json:"organization_id"`
	}
	responseLogedInUser.Username = validatedUser.Username
	responseLogedInUser.Email = validatedUser.Email
	responseLogedInUser.OrganizationID = validatedUser.OrganizationID

	// log.Println("user: %v", validatedUser)
	// if err := DB.First(&validatedUser, validatedUser.ID).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate user"})
	// 	return
	// }
	validatedUser.LastLogin = time.Now()
	DB.Save(&validatedUser)

	c.JSON(http.StatusOK, gin.H{"message": "User is validated", "status:": "success", "user": responseLogedInUser})
}

// Logout logs out a user, if successful, delete cookie header
// Logout godoc
// @Summary Logout a user
// @Description Logout a user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 string {string} json "{"message": "User logged out successfully"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to logout user"}"
// @Router /api/v1/auth/logout [post]
func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}
