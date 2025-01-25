package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// ConnectDB establishes a connection to the MongoDB server
func ConnectDB(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error verifying the MongoDB connection: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB")
	Client = client
	return client
}

// GetCollection returns a reference to a MongoDB collection
func GetCollection(databaseName, collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatalf("MongoDB client is not initialized. Ensure ConnectDB is called before GetCollection.")
	}
	return Client.Database(databaseName).Collection(collectionName)
}

// CloseDB closes the connection to MongoDB
func CloseDB() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := Client.Disconnect(ctx)
		if err != nil {
			log.Fatalf("Error closing the MongoDB connection: %v", err)
		}
		fmt.Println("MongoDB connection closed")
	}
}

// GetMongoURI retrieves the MongoDB URI from environment variables
func GetMongoURI() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}
	return uri
}
