package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI no definido en .env")
	}

	clientOpts := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Error conectando a MongoDB Atlas: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("No se pudo hacer ping a MongoDB: %v", err)
	}

	DB = client.Database("papeleria")
	fmt.Println("✅ Conectado a MongoDB Atlas")
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
