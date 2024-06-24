package controllers

import (
	"github/NjukiG/library-mtaani/initializers"
	"github/NjukiG/library-mtaani/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateAuthor(c *gin.Context) {
	var body struct {
		Name  string
		Email string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	author := models.Author{Name: body.Name, Email: body.Email}

	// Get an admin to create
	user, _ := c.Get("user")
	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to create authors/ Not admin",
		})
		return
	}

	result := initializers.DB.Create(&author)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create a author",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "A new author has been created successfully",
	})
}

func GetAllAuthors(c *gin.Context) {
	var authors []models.Author

	result := initializers.DB.Preload("Books").Find(&authors)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch authors",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Authors": authors,
	})
}

func GetAuthorById(c *gin.Context) {
	authorId := c.Param("id")
	var author models.Author

	result := initializers.DB.Preload("Books").First(&author, authorId)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Author not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Author": author,
	})
}

// I havent written a function to edit the author

// To be done by only an admin
func DeleteAuthor(c *gin.Context) {
	authorId := c.Param("id")

	var author models.Author

	result := initializers.DB.First(&author, authorId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Author not found",
		})
		return
	}

	// Get an admin to delete
	user, _ := c.Get("user")
	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to delete authors/ Not admin",
		})
		return
	}

	initializers.DB.Delete(&author, authorId)

	// Respond
	c.Status(http.StatusNoContent)
	c.JSON(200, gin.H{
		"Message": "An author was deleted...",
	})

}
