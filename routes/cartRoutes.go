package routes

import (
	"github/NjukiG/library-mtaani/controllers"
	"github/NjukiG/library-mtaani/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(r *gin.Engine) {
	protectedRoutes := r.Group("/protected")
	protectedRoutes.Use(middleware.RequireAuth)

	{
		protectedRoutes.POST("/api/cart/:id/items", controllers.AddItemToCart)
		protectedRoutes.DELETE("/api/cart/:id/items/:book_id", controllers.RemoveItemFromCart)
		protectedRoutes.PUT("/api/cart/:id/items/:book_id", controllers.UpdateCartItemQuantity)
		protectedRoutes.GET("/api/cart/:id/items", controllers.ListCartItems)
		protectedRoutes.DELETE("/api/cart/:id/clear", controllers.ClearCart)
		protectedRoutes.GET("/api/cart/:id/review", controllers.ReviewCart)
	}
}
