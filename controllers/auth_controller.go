package controllers

import (
	"net/http"
	"os"
	"parameter-store-be/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user in the database

// Register godoc
// @Summary Register a new user
// @Description Register a new user
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

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newOrganization := model.Organization{
		Name: r.OrganizationName,
	}
	if err := DB.Create(&newOrganization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register organization"})
		return
	}
	organizationID := newOrganization.ID

	newUser := model.User{
		Email:               r.Email,
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
// @Param Cookie header string  false "token"
// @in header
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
	organization := model.Organization{
		Name: l.OrganizationName,
	}

	// Check if the organization exists in the database
	if err := DB.Where("name = ?", l.OrganizationName).First(&organization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	var user model.User
	// Check if the user exists in the database by email and organization_name
	if err := DB.Where("email = ? AND organization_id = ?", l.Email, organization.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	// Generate a JWT token
	jwtToken, err := generateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token for user"})
		return
	}
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		"Authorization",
		jwtToken,
		3600*24*30,
		"",
		"",
		true,
		true,
	)
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "email:": l.Email, "organization_name:": l.OrganizationName, "token": jwtToken})

}

func generateJWTToken(user model.User) (string, error) {
	// Generate a JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenstring, err := jwtToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

// Validate validates a user by cookie
// Validate godoc
// @Summary Validate a user
// @Description Validate a user
// @Tags Auth
// @Accept json
// @Produce json
// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @Success 200 string {string} json "{"message": "User logged in successfully"}"
// @Failure 400 string {string} json "{"error": "Bad request"}"
// @Failure 500 string {string} json "{"error": "Failed to validate user"}"
// @Router /api/v1/auth/validate [get]
func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"message": user})
}
