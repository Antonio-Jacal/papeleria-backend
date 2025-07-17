package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSummary(c *gin.Context) {
	collection := config.GetCollection("pedidos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resultTotal, err := collection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"Hubo un error al encontrar los datos": err})
		return
	}
	var totalWin float64
	var totalPay float64

	for resultTotal.Next(ctx) {
		var total models.List
		if err := resultTotal.Decode(&total); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hubo un error al decodificar los datos": err})
			return
		}
		totalWin += total.TotalGeneral
		totalPay += total.TotalPagado
	}

	result, err := collection.Find(ctx, bson.M{"estaPagado": false})

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"Hubo un error al encontrar los datos": err})
		return
	}

	defer result.Close(ctx)

	var forPayLists []models.List
	var totalForPay float64

	for result.Next(ctx) {
		var lista models.List
		if err := result.Decode(&lista); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hubo un error al decodificar los datos": err})
		}
		forPayLists = append(forPayLists, lista)
		totalForPay += lista.TotalRestante
	}

	c.JSON(http.StatusOK, gin.H{
		"totalGanado":          totalWin,
		"totalPagado":          totalPay,
		"totalPorPagar":        totalForPay,
		"pedidosFaltantesPago": forPayLists,
	})

}

func GetSummaryLabels(c *gin.Context) {

	collection := config.GetCollection("pedidos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.Find(ctx, bson.M{"listaForrada": true}, options.Find().SetProjection(getProjection()))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var labels []models.LabelResumen

	for result.Next(ctx) {
		var labelLists models.LabelResumen
		err = result.Decode(&labelLists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		labels = append(labels, labelLists)
	}

	if err := result.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data": gin.H{
			"labels": labels,
		},
		"message": "Resumen de etiquetas generado correctamente",
	})
}

func getProjection() bson.M {
	return bson.M{
		"_id":                  1,
		"numeroLista":          1,
		"nombreAlumno":         1,
		"grado":                1,
		"fechaEntregaEsperada": 1,
		"etiquetasPersonaje":   1,
		"statusEtiquetas":      1,
		"etiquetasGrandes":     1,
		"etiquetasMedianas":    1,
		"etiquetasChicas":      1,
		"encargadoEtiquetasId": 1,
		"statusForrado":        1,
	}
}
