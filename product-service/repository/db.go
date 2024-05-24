package repository

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() (*mongo.Client, error) {
	mongoUri := os.Getenv("MONGO_URI")
	log.Println("MONGO_URI:", mongoUri)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		return client, err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return client, err
	}

	return client, nil
}
