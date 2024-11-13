package controllers

import (
	"github/NjukiG/library-mtaani/initializers"
	"github/NjukiG/library-mtaani/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Func to add an item to your cart
func AddItemToCart(c *gin.Context) {
	// Assuming user is authenticated
	user, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "User not authenticated",
		})
		return
	}

	userID := user.(models.User).ID

	// Parse the request body
	var body struct {
		BookID   uint `json:"BookID"`
		Quantity int  `json:"Quantity"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request data", "Detail": err.Error()})
		return
	}

	// Find or create the cart for this user
	var cart models.Cart
	if err := initializers.DB.FirstOrCreate(&cart, models.Cart{UserID: userID}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find or create a cart", "Detail": err.Error()})
		return
	}

	// Check if the product exists
	var book models.Book
	if err := initializers.DB.First(&book, body.BookID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Book not found"})
		return
	}

	// Check if the cart item already exists for this product
	var cartItem models.CartItem
	if err := initializers.DB.First(&cartItem, "cart_id = ? AND book_id = ?", cart.ID, body.BookID).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			// Item does not exist in cart, create new one
			cartItem = models.CartItem{
				CartID:   cart.ID,
				BookID:   body.BookID,
				Quantity: body.Quantity,
				UserID:   userID, // Make sure UserID is set
			}

			if err := initializers.DB.Create(&cartItem).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to add item to cart", "Detail": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"Message": "Item added to cart"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to check cart items", "Detail": err.Error()})
			return
		}
	}

	// If item already exists, update the quantity
	cartItem.Quantity += body.Quantity
	if err := initializers.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update cart item quantity", "Detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Cart item updated"})
}

// Func to list all items in cart
func ListCartItems(c *gin.Context) {
	userId := c.Param("id")
	// cartID := c.Param("id")

	var cartItems []models.CartItem

	// resultItems := initializers.DB.Where("cart_id = ?", cartID).Preload("Book").Find(&cartItems)
	resultItems := initializers.DB.Where("user_id = ?", userId).Preload("Book").Find(&cartItems)

	if resultItems.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "Cart items not found",
		})
		return
	}
	totalPrice := 0
	for _, item := range cartItems {
		totalPrice += item.Book.Price * item.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"CartItems": cartItems,
		// "Total Price": totalPrice,
	})
}

// Func to remove an item from your cart
func RemoveItemFromCart(c *gin.Context) {
	// cartID := c.Param("id")
	userID := c.Param("id")
	bookID := c.Param("book_id")

	var cartItem models.CartItem

	resultItem := initializers.DB.Where("user_id = ? AND book_id = ?", userID, bookID).First(&cartItem)

	if resultItem.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cart item not found",
		})
		return
	}

	result := initializers.DB.Delete(&cartItem)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to remove item from cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart",
	})
}

// Fuc to update the quantity of your cart items
func UpdateCartItemQuantity(c *gin.Context) {
	cartID := c.Param("id")
	bookID := c.Param("book_id")

	var body struct {
		Quantity int `json:"quantity"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	var cartItem models.CartItem
	if err := initializers.DB.Where("cart_id = ? AND book_id = ?", cartID, bookID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cart item not found",
		})
		return
	}

	// Update the quantity
	cartItem.Quantity = body.Quantity
	if err := initializers.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update cart item quantity",
		})
		return
	}

	// Preload the associated Product data
	initializers.DB.Preload("Book").First(&cartItem, cartItem.ID)

	c.JSON(http.StatusOK, gin.H{
		"CartItem": cartItem,
	})
}

// Func to clear the whole cart
func ClearCart(c *gin.Context) {
	userID := c.Param("id")
	// cartID := c.Param("id")

	result := initializers.DB.Where("user_id = ?", userID).Delete(&models.CartItem{})

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to clear cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Cart cleared",
	})
}

// Func get all the details, price, quantities e.t.c
func ReviewCart(c *gin.Context) {
	userID := c.Param("id")
	// cartID := c.Param("id")

	var cartItems []models.CartItem

	resultItems := initializers.DB.Where("user_id = ?", userID).Preload("Book").Find(&cartItems)

	if resultItems.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "Cart items not found",
		})
		return
	}
	totalPrice := 0
	for _, item := range cartItems {
		totalPrice += item.Book.Price * item.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"CartItems":  cartItems,
		"TotalPrice": totalPrice,
	})
}
