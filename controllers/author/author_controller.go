// controllers/author/author_controller.go

package author

import (
	"net/http"
	"strconv"
	main "vcs_backend/gorm/controllers"
	"vcs_backend/gorm/models"

	"github.com/gin-gonic/gin"
)

// CREATE authors
func RegisterAuthor(c *gin.Context) {
	var registingAuthor models.Author
	registingAuthor.FirstName = c.Query("first-name")
	registingAuthor.LastName = c.Query("last-name")
	registingAuthor.Email = c.Query("email")
	registingAuthor.Phone = c.Query("phone")
	registingAuthor.Address = c.Query("adress")
	registingAuthor.Password = c.Query("password")

	// Bind the JSON body to the registingAuthor struct
	// if err := c.BindJSON(&registingAuthor); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error BindJSON": err.Error()})
	// 	return
	// }

	// Create the registingAuthor in the database
	if err := main.DB.Create(&registingAuthor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Create": "Failed to register author"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Author registered successfully", "author": registingAuthor})
}

// READ authors
func GetAuthors(c *gin.Context) {
	var authors []models.Author
	main.DB.Find(&authors)
	// fmt.Println(authors)
	c.JSON(http.StatusOK, gin.H{
		"authors": authors,
	})
}
func GetAuthorById(c *gin.Context) {
	var author []models.Author
	idStr := c.Query("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}
	main.DB.Find(&author, id)

	c.JSON(http.StatusOK, gin.H{
		"authors": author,
	})
}

func GetAuthorsByName(c *gin.Context) {
	var authors []models.Author
	firstName := c.Query("first-name")
	lastName := c.Query("last-name")

	// Construct the conditional query string based on the provided first name and last name
	var conditionalString string
	var conditions []interface{}
	if firstName != "" && lastName != "" {
		conditionalString = "first_name = ? AND last_name = ?"
		conditions = append(conditions, firstName, lastName)
	} else if firstName != "" {
		conditionalString = "first_name = ?"
		conditions = append(conditions, firstName)
	} else if lastName != "" {
		conditionalString = "last_name = ?"
		conditions = append(conditions, lastName)
	}

	// Execute the query with the constructed conditional string and conditions
	if err := main.DB.Where(conditionalString, conditions...).Find(&authors).Error; err != nil {
		panic("Error when GetAuthorsByName")
	}

	c.JSON(http.StatusOK, gin.H{
		"authors": authors,
	})
}

// UPDATE authors
func UpdateAuthorInfo(c *gin.Context) {
	var author models.Author
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	// Find the author by ID
	if err := main.DB.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	// Bind the JSON body to update author information
	if err := c.BindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated author back to the database
	if err := main.DB.Save(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Author updated successfully", "author": author})
}
func DeleteAuthor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}
	var author models.Author
	if err := main.DB.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
	}
	if err := main.DB.Delete(&author).Error; err != nil {
		// Display JSON error
		c.JSON(404, gin.H{"error": "User not found"})
	}
}
