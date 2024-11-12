package routes

import (
	"github/NjukiG/library-mtaani/controllers"
	"github/NjukiG/library-mtaani/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(r *gin.Engine) {

	publicRoutes := r.Group("/public")

	{
		publicRoutes.GET("/api/books", controllers.GetAllBooks)
		publicRoutes.GET("/api/books/:id", controllers.GetBookById)
		publicRoutes.GET("/api/categories/:id/books", controllers.GetBooksByCategory)

	}

	protectedRoutes := r.Group("/protected")
	protectedRoutes.Use(middleware.RequireAuth)

	{
		protectedRoutes.POST("/api/books", controllers.PostNewBook)
		protectedRoutes.PUT("/api/books/:id", controllers.EditBookDetails)
		protectedRoutes.DELETE("/api/books/:id", controllers.DeleteBook)
	}

}
