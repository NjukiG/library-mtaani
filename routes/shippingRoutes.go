package routes

import (
	"github/NjukiG/library-mtaani/controllers"
	"github/NjukiG/library-mtaani/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterShippingRoutes(r *gin.Engine) {
	protectedRoutes := r.Group("/protected")
	protectedRoutes.Use(middleware.RequireAuth)

	{
		protectedRoutes.POST("/api/shipping-details", controllers.AddShippingDetails)
		protectedRoutes.PUT("/api/shipping-details", controllers.UpdateShippingDetails)
		protectedRoutes.GET("/api/shipping-details", controllers.GetShippingDetails)
	}
}
