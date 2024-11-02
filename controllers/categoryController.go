package controllers

import (
	"github/NjukiG/library-mtaani/initializers"
	"github/NjukiG/library-mtaani/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddNewCategory(c *gin.Context) {

	var body struct {
		Title string
		Books []models.Book
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user, _ := c.Get("user")

	category := models.Category{
		Title: body.Title,
	}

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to post a book/ Not an admin",
		})
		return
	}

	result := initializers.DB.Create(&category)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create the category...",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Category": category,
	})
}


func GetAllCategories(c *gin.Context){
	var categories []models.Category

	result := initializers.DB.Preload("Books").Find(&categories)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Categories": categories,
	})
}


func GetCategoryById (c *gin.Context){
	id := c.Param("id")

	var category models.Category

	result := initializers.DB.Preload("Books").First(&category, id)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Category": category,
	})
}
