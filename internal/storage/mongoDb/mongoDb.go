package mongoDb

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
	err := godotenv.Load("internal/config/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	mongoDbUri := "mongodb:" + os.Getenv("DATABASE_HOST") + ":" + os.Getenv("DATABASE_PORT")
	if os.Getenv("DATABASE_HOST") == "" || os.Getenv("DATABASE_PORT") == "" {
		log.Fatal("Error connecting to MongoDB")
	}

	clientOptions := options.Client().ApplyURI(mongoDbUri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}

	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	err := godotenv.Load("internal/config/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")

	collection := client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return nil
	}

	return collection
}
