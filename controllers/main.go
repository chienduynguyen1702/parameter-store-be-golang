package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"parameter-store-be/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// SetDB sets the db object
func SetDB(database *gorm.DB) {
	DB = database
}

// generateBcryptPassword generates a bcrypt hash of the password
func generateBcryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// generateJWTToken generates a JWT token
func generateJWTToken(user models.User) (string, error) {
	// Generate a JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"org_id":  user.OrganizationID,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenstring, err := jwtToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

// gen Token for agents, by agentID and orgID
func GenerateTokenForAgent(agentID, orgID string) string {
	plainText := agentID + orgID
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash to store:", string(hash))

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

// func parseJWTTokenFromCookie(c *gin.Context) (jwt.MapClaims, error) {
// 	tokenString, err := c.Cookie("Authorization")
// 	if err != nil {
// 		return nil, err
// 	}
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Don't forget to validate the alg is what you expect:
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, jwt.ErrSignatureInvalid
// 		}
// 		return []byte(os.Getenv("SECRET_KEY")), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		return claims, nil
// 	}
// 	return nil, jwt.ErrSignatureInvalid
// }

// getUserFromContext retrieves the user from the context
func getUserFromContext(c *gin.Context) (models.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return models.User{}, errors.New("Failed to get user from context")
	}
	return user.(models.User), nil
}

func paginationDataParam(paramList []models.Parameter, page, limit int) []models.Parameter {
	start := (page - 1) * limit
	end := page * limit
	if start > len(paramList) {
		start = len(paramList)
	}
	if end > len(paramList) {
		end = len(paramList)
	}
	return paramList[start:end]
}
func paginationDataAgent(agentList []models.Agent, page, limit int) []models.Agent {
	start := (page - 1) * limit
	end := page * limit
	if start > len(agentList) {
		start = len(agentList)
	}
	if end > len(agentList) {
		end = len(agentList)
	}
	return agentList[start:end]
}

func paginationDataUser(userList []models.User, page, limit int) []models.User {
	start := (page - 1) * limit
	end := page * limit
	if start > len(userList) {
		start = len(userList)
	}
	if end > len(userList) {
		end = len(userList)
	}
	return userList[start:end]
}

type DashboardCard struct {
	Icon    string `json:"icon"`
	Value   int    `json:"value"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Tooltip string `json:"tooltip"`
}
