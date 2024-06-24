package controllers

import (
	"github/NjukiG/library-mtaani/initializers"
	"github/NjukiG/library-mtaani/models"
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
)

// New attempt to borrow book
func BorrowBook(c *gin.Context) {
	bookId := c.Param("id")
	var book models.Book

	result := initializers.DB.Preload("Borrows").First(&book, bookId)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Book not found",
		})
		return
	}

	var body struct {
		DueDate string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user, _ := c.Get("user")

	borrow := models.Borrow{
		UserID:  user.(models.User).ID,
		BookID:  book.ID,
		DueDate: body.DueDate,
		// Book: book,
		// User: user,
	}

	newBorrow := initializers.DB.Create(&borrow)

	if newBorrow.Error != nil {
		c.Status(400)
		return
	}

	book.Copies -= 1
	result2 := initializers.DB.Save(&book)
	if result2.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update book copies",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book borrowed successfully"})

}

func GetBorrowedBooks(c *gin.Context) {
	userID, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var borrows []models.Borrow
	result := initializers.DB.Preload("Book").Where("user_id = ?", userID.(models.User).ID).Find(&borrows)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch borrowed books",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"borrowed_books": borrows})
}

func ReturnBorrowedBooks(c *gin.Context) {
	userID, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var body struct {
		BorrowID uint
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var borrow models.Borrow
	if err := initializers.DB.First(&borrow, body.BorrowID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Borrow record not found"})
		return
	}

	if borrow.UserID != userID.(models.User).ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can only return your own borrowed books"})
		return
	}

	var book models.Book
	if err := initializers.DB.First(&book, borrow.BookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err := initializers.DB.Delete(&borrow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to return book"})
		return
	}

	book.Copies += 1
	if result := initializers.DB.Save(&book); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book copies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}
