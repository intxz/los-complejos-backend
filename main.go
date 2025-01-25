// main.go
package main

import (
	"log"
	"los-complejos-backend/database"
	"los-complejos-backend/handlers"
	"los-complejos-backend/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Message struct for test endpoint response
type Message struct {
	Content string `json:"content"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set in the environment")
	}

	uri := "mongodb://localhost:27017"

	// Connect to the database
	_ = database.ConnectDB(uri)
	defer database.CloseDB()

	// Collections
	complejo_collection := database.GetCollection("COMPLEJOS", "complejo")
	event_collection := database.GetCollection("COMPLEJOS", "event")

	r := gin.Default()

	// Test route
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, Message{Content: "Server is running!"})
	})

	// Complejo routes
	// Handles user management for "Complejo" resources
	r.POST("/complejo", handlers.CreateComplejo(complejo_collection))
	r.GET("/complejo", handlers.GetComplejos(complejo_collection))
	r.GET("/complejo/:id", handlers.GetComplejo(complejo_collection))
	r.PUT("/complejo/admin", middleware.AuthMiddleware(), handlers.UpdateComplejoForAdmin(complejo_collection))
	r.PUT("/complejo/user", middleware.AuthMiddleware(), handlers.UpdateComplejoForUser(complejo_collection))

	// Event routes
	// Handles event management and user subscription/unsubscription
	r.POST("/event", middleware.AuthMiddleware(), handlers.CreateEvent(event_collection))
	r.GET("/event", handlers.GetEvents(event_collection))
	r.GET("/event/:id", handlers.GetEvent(event_collection))
	r.PUT("/event/admin", middleware.AuthMiddleware(), handlers.UpdateEventForAdmin(event_collection))
	r.PUT("/event/:id/subscribe", middleware.AuthMiddleware(), handlers.SubscribeEvent(event_collection))
	r.PUT("/event/:id/unsubscribe", middleware.AuthMiddleware(), handlers.UnsuscribeEvent(event_collection))

	// Start the server on port 8080
	r.Run(":8080")
}
