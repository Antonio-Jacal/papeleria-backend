package utils

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenerateNextNumeroLista(collection *mongo.Collection) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.FindOne().SetSort(bson.D{{Key: "numeroLista", Value: -1}})
	var lastList models.List

	err := collection.FindOne(ctx, bson.M{
		"numeroLista": bson.M{"$regex": "^LTS\\d+$"},
	}, opts).Decode(&lastList)

	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}

	nextNumber := 1
	if lastList.NumeroLista != "" {
		fmt.Sscanf(lastList.NumeroLista, "LTS%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("LTS%d", nextNumber), nil
}

func GeneratePin() string {
	min := 1000
	max := 9999
	rand.Seed(time.Now().UnixNano())
	numero := rand.Intn(max-min+1) + min
	return fmt.Sprintf("%d", numero)
}
