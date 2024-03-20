package controllers

import (
	"net/http"
	"parameter-store-be/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CREATE authors
func RegisterAuthor(c *gin.Context) {
	var registingAuthor models.Author
	registingAuthor.FirstName = c.Query("first-name")
	registingAuthor.LastName = c.Query("last-name")
	registingAuthor.Email = c.Query("email")
	registingAuthor.Phone = c.Query("phone")
	registingAuthor.Address = c.Query("address")
	registingAuthor.Password = c.Query("password")

	// Bind the JSON body to the registingAuthor struct
	// if err := c.BindJSON(&registingAuthor); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error BindJSON": err.Error()})
	// 	return
	// }

	// Create the registingAuthor in the database
	if err := DB.Create(&registingAuthor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Create": "Failed to register author"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Author registered successfully", "author": registingAuthor})
}

// CREATE new post
func CreateNewPost(c *gin.Context) {
	authorID := c.Param("id") // Get author ID from the URL parameter
	var author models.Author

	// Check if the author exists
	if err := DB.First(&author, authorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	var newPost models.Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	// Create the post
	if err := DB.Create(&newPost).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Associate the post with the author
	author.Posts = append(author.Posts, newPost)
	if err := DB.Save(&author).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    newPost,
	})
}

// READ authors
func GetAllAuthors(c *gin.Context) {
	var authors []models.Author
	DB.Find(&authors)
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
	DB.Find(&author, id)

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
	if err := DB.Where(conditionalString, conditions...).Find(&authors).Error; err != nil {
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
	if err := DB.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	// Bind the JSON body to update author information
	if err := c.BindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated author back to the database
	if err := DB.Save(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Author updated successfully", "author": author})
}

// DELETE
func DeleteAuthor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}
	var author models.Author
	if err := DB.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	} else if err := DB.Delete(&author).Error; err != nil {
		// Display JSON error
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
}
func DeletePostOfAuthor(c *gin.Context) {

	// Parse author ID from URL parameter
	authorID, err := strconv.Atoi(c.Param("author-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid author ID",
		})
		return
	}

	// Retrieve author based on author ID
	var author models.Author
	if err := DB.First(&author, authorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Author not found",
		})
		return
	}

	// Parse post ID from query parameter
	postID, err := strconv.Atoi(c.Query("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	// Check if the post exists
	var post models.Post
	if err := DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Post not found",
		})
		return
	}
	// Remove the association between the author and the post
	if err := DB.Table("author_posts").Where("author_id = ? AND post_id = ?", authorID, postID).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete association",
		})
		return
	}
	// Delete post
	if err = DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
