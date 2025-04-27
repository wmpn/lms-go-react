package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection

func ConnectMongoDB() error {
	// var err error
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(MONGO_URI)


	for attempt := 1; attempt <= 3; attempt++ {
		// Create a context that automatically times out after 10 seconds
		// if the MongoDB server is slow or unreachable
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel() // Ensure the short-lived context is cancelled to avoid resource leaks

		client, connErr := mongo.Connect(ctx, clientOptions)
		if connErr != nil {
			err = connErr
			log.Printf("[Attempt %d] Failed to connect to MongoDB: %v", attempt, connErr)
			continue
		}

		// Ping the database to ensure the connection is alive
		if pingErr := client.Ping(ctx, nil); pingErr != nil {
			err = pingErr
			log.Printf("[Attempt %d] Failed to ping MongoDB: %v", attempt, pingErr)
			continue
		}

		fmt.Println("✅ Connected to MongoDB successfully!")
		Client = client
		Collection = Client.Database("lms_db").Collection("books")
		return nil // return nil as the error is nil in return of the function
	}
	
	// Create a new formatted error message with the original error wrapped inside
	// using the %w verb to indicate that the error is wrapped.
	// returned error: "failed to connect to MongoDB: connection refused"
	return fmt.Errorf("failed to connect to MongoDB after retries: %w", err)
}
