package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/wmpn/lms-go-react/db"
	"github.com/wmpn/lms-go-react/routes"
)

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	// Connect to MongoDB
	if err := db.ConnectMongoDB(); err != nil {
		log.Fatal("MongoDB connection failed: ", err)
	}
	
	// Setup routes
	routes.RegisterBookRoutes(app)
	
	// Start the server
	// Listen on the specified port
	log.Fatal(app.Listen(":" + PORT))
}
