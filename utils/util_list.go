package utils

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/gin-gonic/gin"
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

func PrintAllQueryParams(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	if len(queryParams) == 0 {
		fmt.Println("No se recibieron parámetros en la query.")
		return
	}

	fmt.Println("Parámetros recibidos:")
	for key, values := range queryParams {
		// En caso de que haya múltiples valores para una misma clave
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
}

// ParseTimeParam convierte un string a time.Time
func ParseTimeParam(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	// Primero intenta con formato simple YYYY-MM-DD
	t, err := time.Parse("2006-01-02", dateStr)
	if err == nil {
		return &t, nil
	}

	// Si falla, intenta con formato RFC3339
	t, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("fecha debe estar en formato YYYY-MM-DD o RFC3339")
	}

	return &t, nil
}
