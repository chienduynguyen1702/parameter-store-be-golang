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
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
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
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check if the token is expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find the user in the database
		var user models.User
		controllers.DB.First(&user, claims["user_id"])
		log.Printf("User found in the database %v with ID %v", user, claims["user_id"])

		if user.ID == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
			return
		}
		// Set the user and their organization_id in the context
		c.Set("user", user)
		orgID, ok := claims["org_id"].(float64)
		if !ok {
			log.Fatal("Failed to parse org_id as float64")
		}

		// Convert float64 to uint
		orgIDUint := uint(orgID)

		// Set org_id in the context
		c.Set("org_id", orgIDUint)
	}
	c.Next()
}
