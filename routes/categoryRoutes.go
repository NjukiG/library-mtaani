package routes

import (
	"github/NjukiG/library-mtaani/controllers"
	"github/NjukiG/library-mtaani/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(r *gin.Engine) {

	publicRoutes := r.Group("/public")
	{
		publicRoutes.GET("/api/categories", controllers.GetAllCategories)
		publicRoutes.GET("/api/categories/:id", controllers.GetCategoryById)
	}

	protectedRoutes := r.Group("/protected")
	protectedRoutes.Use(middleware.RequireAuth)

	{
		protectedRoutes.POST("/api/categories", controllers.AddNewCategory)
	}

}
