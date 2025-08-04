package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EncargadoEtiquetasRequest struct {
	EncargadoEtiquetas string `json:"encargadoEtiquetasId"`
	StatusEtiquetas    string `json:"statusEtiquetas"`
}

type requestPedidoUpdate struct {
	EstadoLista  string         `json:"estado_lista,omitempty"`
	Faltantes    map[string]int `json:"faltantes,omitempty"`
	EstaPagado   bool           `json:"estaPagado,omitempty"`
	Pago         models.Pago    `json:"pago,omitempty"`
	TotalGeneral float64        `json:"totalGeneral,omitempty"`
}

func UpdateLabelResponse(c *gin.Context) {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var responsable EncargadoEtiquetasRequest

	if err := c.ShouldBindJSON(&responsable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"encargadoEtiquetasId": responsable.EncargadoEtiquetas,
			"statusEtiquetas":      responsable.StatusEtiquetas,
		},
	}

	collection := config.GetCollection("pedidos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.UpdateByID(ctx, objectID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "Usuario actualizado correctamente"})

}

func UpdatePedido(c *gin.Context) {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var request requestPedidoUpdate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	setFields := bson.M{}
	pushFields := bson.M{}

	if request.EstadoLista != "" {
		setFields["estadoLista"] = request.EstadoLista
	}

	if request.Faltantes != nil {
		setFields["faltantes"] = request.Faltantes
	}

	if request.EstaPagado {
		setFields["estaPagado"] = true
		setFields["totalRestante"] = 0
		setFields["totalPagado"] = request.TotalGeneral
		pushFields["pagos"] = request.Pago
	}

	update := bson.M{}
	if len(setFields) > 0 {
		update["$set"] = setFields
	}
	if len(pushFields) > 0 {
		update["$push"] = pushFields
	}
	collection := config.GetCollection("pedidos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Update query:", update)

	result, err := collection.UpdateByID(ctx, objectID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "lista no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Usuario actualizado correctamente",
		"status":  "ok",
	})
}

func UpdateForrado(c *gin.Context) {

	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	type requestUpdate struct {
		StatusForradoUtiles string `json:"statusForradoUtiles,omitempty"`
		StatusForradoLibros string `json:"statusForradoLibros,omitempty"`
	}

	var request requestUpdate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if request.StatusForradoLibros == "" {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "statusForradoLibros is empty"})
		return
	}
	if request.StatusForradoUtiles == "" {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": "statusForradoUtiles is empty"})
		return
	}

	setFields := bson.M{}

	setFields["statusForradoUtiles"] = request.StatusForradoUtiles
	setFields["statusForradoLibros"] = request.StatusForradoLibros

	if request.StatusForradoLibros == "Forrados" && request.StatusForradoUtiles == "Etiquetados" {
		setFields["statusForrado"] = "Forrada"
	}

	collection := config.GetCollection("pedidos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println(setFields)

	_, err = collection.UpdateByID(ctx, objectID, bson.M{"$set": setFields})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ocurrio un error al actulizar: ": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Se actualizo el forrado correctamente"})

}
