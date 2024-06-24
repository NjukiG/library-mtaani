package controllers

import (
	"github/NjukiG/library-mtaani/initializers"
	"github/NjukiG/library-mtaani/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostNewBook(c *gin.Context) {
	authorId := c.Param("id")
	var author models.Author

	myAuthor := initializers.DB.Preload("Books").First(&author, authorId)

	if myAuthor.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Author not found",
		})
		return
	}

	var body struct {
		Title       string
		ImageUrl    string
		Price       float64
		Copies      int64
		Description string
		// AuthorID    uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	user, _ := c.Get("user")
	book := models.Book{
		Title:       body.Title,
		ImageUrl:    body.ImageUrl,
		Price:       body.Price,
		Copies:      body.Copies,
		Description: body.Description,
		AuthorID:    author.ID,
		// Author:      author.Name,
	}

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to post a book/ Not an admin",
		})
		return
	}

	result := initializers.DB.Create(&book)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create the book...",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Book": book,
	})
}

func GetAllBooks(c *gin.Context) {
	var books []models.Book

	result := initializers.DB.Preload("Borrows").Find(&books)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Books not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Books": books,
	})
}

func GetBookById(c *gin.Context) {
	bookId := c.Param("id")

	var book models.Book

	result := initializers.DB.Preload("Borrows").First(&book, bookId)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Book with that id not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Book": book,
	})
}

func EditBookDetails(c *gin.Context) {
	bookId := c.Param("id")

	var body struct {
		Title       string
		ImageUrl    string
		Price       float64
		Copies      int64
		Description string
		AuthorID    uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	var book models.Book

	initializers.DB.First(&book, bookId)

	user, _ := c.Get("user")

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to edit book details",
		})
		return
	}

	initializers.DB.Model(&book).Updates(models.Book{
		Title:       body.Title,
		ImageUrl:    body.ImageUrl,
		Price:       body.Price,
		Copies:      body.Copies,
		Description: body.Description,
		// AuthorID:    body.AuthorID,
	})

	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	bookId := c.Param("id")

	var book models.Book

	result := initializers.DB.First(&book, bookId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Book not found",
		})
		return
	}

	user, _ := c.Get("user")

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to delete a book / Not an admin",
		})
		return
	}

	initializers.DB.Delete(&book, bookId)

	// Respond
	c.Status(http.StatusNoContent)
	c.JSON(200, gin.H{
		"Message": "A book was deleted...",
	})

}
