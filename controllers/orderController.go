package controllers

import (
	"github/NjukiG/library-mtaani/initializers"
	"github/NjukiG/library-mtaani/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateOrder creates a new order based on the cart items
func CreateOrder(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(models.User).ID

	// Find the cart for this user
	var cart models.Cart
	initializers.DB.Preload("CartItems.Book").Where("user_id = ?", userID).First(&cart)

	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your cart is empty!"})
		return
	}

	// Find the shipping details for this user
	var shippingDetail models.ShippingDetail
	initializers.DB.Where("user_id = ?", userID).First(&shippingDetail)

	// Calculate total price
	totalPrice := 0
	orderItems := []models.OrderItem{}
	for _, cartItem := range cart.CartItems {
		totalPrice += cartItem.Book.Price * cartItem.Quantity
		orderItems = append(orderItems, models.OrderItem{
			BookID:   cartItem.BookID,
			Quantity: cartItem.Quantity,
			Price:    cartItem.Book.Price,
		})
	}

	// Create the order
	order := models.Order{
		UserID:           userID,
		TotalPrice:       totalPrice,
		ShippingDetailID: shippingDetail.ID,
		OrderItems:       orderItems,
	}

	if err := initializers.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create order"})
		return
	}

	// Clear the cart
	initializers.DB.Where("user_id = ?", userID).Delete(&models.CartItem{})

	c.JSON(http.StatusCreated, gin.H{"Message": "Order created successfully", "Order": order})
}

// Get all orders for a user
func GetOrders(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(models.User).ID

	var orders []models.Order
	initializers.DB.Preload("OrderItems.Book").Where("user_id = ?", userID).Find(&orders)

	c.JSON(http.StatusOK, gin.H{"Orders": orders})
}

// Get single order by ID
func GetOrder(c *gin.Context) {
	// Extract the order ID from the URL parameters
	user, _ := c.Get("user")
	userID := user.(models.User).ID
	var orders []models.Order
	result := initializers.DB.Preload("OrderItems.Book").Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve orders",
		})
		return
	}
	// Extract the order ID from the URL parameters and convert it to a uint
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	// Loop through the user's orders to find the specific order by ID
	for _, order := range orders {
		if order.ID == uint(orderID) {
			c.JSON(http.StatusOK, gin.H{
				"order": order,
			})
			return
		}
	}

	// If no matching order is found, return a 404 error
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Order not found",
	})
}

// Update order status
// Update Order Status handles the request to update an order's status
func UpdateOrderStatus(c *gin.Context) {
	// Fetch the order using the order_id from the URL parameters
	var order models.Order
	if err := initializers.DB.Where("id = ?", c.Param("order_id")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Bind the request body to get the new status
	var body struct {
		Status models.OrderStatus `json:"status" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Validate the new status
	validStatuses := []models.OrderStatus{
		models.OrderPending,
		models.OrderProcessed,
		models.OrderShipped,
		models.OrderDelivered,
		models.OrderCancelled,
	}

	isValidStatus := false
	for _, status := range validStatuses {
		if body.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
		return
	}

	// Update the order status
	order.Status = string(body.Status)
	if err := initializers.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully", "order": order})
}
