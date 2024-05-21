package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"parameter-store-be/controllers"
	"parameter-store-be/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequiredAuth(c *gin.Context) {
	// tokenString, err := c.Cookie("Authorization")
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get token"})
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }
	tokenString := c.GetHeader("Authorization")
	// log.Printf("debug: tokenString \"%s\"", tokenString)
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to get token in header"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check if the token is expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
			return
		}

		// Find the user in the database
		var user models.User
		controllers.DB.First(&user, claims["user_id"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		if user.IsArchived {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User is archived"})
			return
		}
		// Set the user and their organization_id in the context
		c.Set("user", user)
		orgID, ok := claims["org_id"].(float64)
		if !ok {
			log.Fatal("Failed to parse org_id as float64")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse org_id as float64"})
			return
		}

		// Convert float64 to uint
		orgIDUint := uint(orgID)

		// Set org_id in the context
		c.Set("org_id", orgIDUint)
	}
	c.Next()
}
