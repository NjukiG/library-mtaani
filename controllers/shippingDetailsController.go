package controllers

import (
	"github/NjukiG/library-mtaani/initializers"
	"github/NjukiG/library-mtaani/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddShippingDetails(c *gin.Context) {
	var body struct {
		Address     string
		City        string
		State       string
		PostalCode  string
		Country     string
		PhoneNumber string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to read request body"})
		return
	}

	// Get the authenticated user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "User not authenticated"})
		return
	}
	userID := user.(models.User).ID

	// Check if user already has shipping details
	var existingShippingDetail models.ShippingDetail
	if err := initializers.DB.Where("user_id = ?", userID).First(&existingShippingDetail).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"Error": "User already has shipping details",
		})
		return
	}

	// Add new shipping details
	shippingDetail := models.ShippingDetail{
		UserID:      userID,
		Address:     body.Address,
		City:        body.City,
		State:       body.State,
		PostalCode:  body.PostalCode,
		Country:     body.Country,
		PhoneNumber: body.PhoneNumber,
	}

	if result := initializers.DB.Create(&shippingDetail); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to add shipping details"})
		return
	}

	// Reload the shipping detail with user data to return the full object
	if err := initializers.DB.Preload("User").First(&shippingDetail, shippingDetail.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to load shipping details with user data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ShippingDetails": shippingDetail})
}

// UpdateShippingDetails updates existing shipping details for the user
func UpdateShippingDetails(c *gin.Context) {
	var body struct {
		Address     string
		City        string
		State       string
		PostalCode  string
		Country     string
		PhoneNumber string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Get the user from the context
	user, _ := c.Get("user")
	userID := user.(models.User).ID

	var shippingDetail models.ShippingDetail

	// Find the shipping detail by user ID
	if err := initializers.DB.Where("user_id = ?", userID).First(&shippingDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping details not found"})
		return
	}

	// Update the shipping details
	shippingDetail.Address = body.Address
	shippingDetail.City = body.City
	shippingDetail.State = body.State
	shippingDetail.PostalCode = body.PostalCode
	shippingDetail.Country = body.Country
	shippingDetail.PhoneNumber = body.PhoneNumber

	if err := initializers.DB.Save(&shippingDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shipping details"})
		return
	}

	// Reload the shipping detail with user data to return the full object
	if err := initializers.DB.Preload("User").First(&shippingDetail, shippingDetail.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to load shipping details with user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message":         "Shipping details updated successfully",
		"ShippingDetails": shippingDetail,
	})
}

func GetShippingDetails(c *gin.Context) {
	// Get the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "User not authenticated"})
		return
	}
	userID := user.(models.User).ID

	var shippingDetail models.ShippingDetail

	// Preload the User, Cart, and CartItems associations
	if err := initializers.DB.Preload("User.Cart.CartItems.Book").First(&shippingDetail, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping details not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ShippingDetails": shippingDetail})
}
