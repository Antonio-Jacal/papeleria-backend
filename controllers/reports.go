package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetSummary(c *gin.Context) {
	collection := config.GetCollection("pedidos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.Find(ctx, bson.M{"estaPagado": false})

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"Hubo un error al encontrar los datos": err})
		return
	}

	defer result.Close(ctx)

	var forPayLists []models.List
	var totalWin float64
	var totalPay float64
	var totalForPay float64

	for result.Next(ctx) {
		var lista models.List
		if err := result.Decode(&lista); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hubo un error al decodificar los datos": err})
		}
		forPayLists = append(forPayLists, lista)
		totalWin += lista.TotalGeneral
		totalPay += lista.TotalPagado
		totalForPay += lista.TotalRestante
	}

	c.JSON(http.StatusOK, gin.H{
		"totalGanado":          totalWin,
		"totalPagado":          totalPay,
		"totalPorPagar":        totalForPay,
		"pedidosFaltantesPago": forPayLists,
	})
}
