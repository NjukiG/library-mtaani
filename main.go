package main

import (
	"fmt"
	"github/NjukiG/library-mtaani/initializers"
	"github/NjukiG/library-mtaani/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	fmt.Println("Library Mtaani System...")

	r := gin.Default()

	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Adjust this as needed
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterRoutes(r)
	routes.RegisterCategoryRoutes(r)
	routes.RegisterBookRoutes(r)
	routes.RegisterCartRoutes(r)
	routes.RegisterShippingRoutes(r)
	routes.RegisterOrderRoutes(r)

	r.Run()
}
