package routes

import (
	"github/NjukiG/library-mtaani/controllers"
	"github/NjukiG/library-mtaani/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	publicRoutes := r.Group("/public")
	{
		publicRoutes.POST("/api/register", controllers.Register)
		publicRoutes.POST("/api/login", controllers.Login)
	}

	protectedRoutes := r.Group("/protected")
	protectedRoutes.Use(middleware.RequireAuth)

	{
		protectedRoutes.GET("/api/validate", controllers.Validate)
		protectedRoutes.POST("/api/logout", controllers.Logout)

		// Author routes
		protectedRoutes.POST("/api/authors", controllers.CreateAuthor)
		protectedRoutes.GET("/api/authors", controllers.GetAllAuthors)
		protectedRoutes.GET("/api/authors/:id", controllers.GetAuthorById)
		protectedRoutes.DELETE("/api/authors/:id", controllers.DeleteAuthor)

		// Borrow books routes
		protectedRoutes.POST("/api/books/:id/borrows", controllers.BorrowBook)
		protectedRoutes.GET("/api/borrows", controllers.GetBorrowedBooks)
		protectedRoutes.POST("/api/borrows/return", controllers.ReturnBorrowedBooks)

	}
}
