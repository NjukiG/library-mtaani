package main

import (
	"fmt"
	"github/NjukiG/library-mtaani/initializers"
	"github/NjukiG/library-mtaani/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	fmt.Println("Library Mtaani System...")

	r := gin.New()
	routes.RegisterRoutes(r)

	// publicRoutes := r.Group("/public")
	// {
	// 	publicRoutes.POST("/api/register", controllers.Register)
	// 	publicRoutes.POST("/api/login", controllers.Login)
	// }

	// protectedRoutes := r.Group("/protected")
	// protectedRoutes.Use(middleware.RequireAuth)

	// {
	// 	protectedRoutes.GET("/api/validate", controllers.Validate)
	// 	protectedRoutes.POST("/api/logout", controllers.Logout)
	// }

	r.Run()
}
