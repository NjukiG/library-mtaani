package routes

import (
	"github/NjukiG/library-mtaani/controllers"
	"github/NjukiG/library-mtaani/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(router *gin.Engine) {
	protectedRoutes := router.Group("/protected")
	protectedRoutes.Use(middleware.RequireAuth)

	// Order routes
	{
		protectedRoutes.POST("/api/cart/:id/orders", controllers.CreateOrder)
		protectedRoutes.GET("/api/cart/:id/orders", controllers.GetOrders)
		protectedRoutes.GET("/api/cart/:id/orders/:id", controllers.GetOrder)
		// protectedRoutes.PUT("/orders/:order_id", controllers.UpdateOrderStatus)
		// protectedRoutes.DELETE("/orders/:order_id", controllers.DeleteOrder)
	}
}
