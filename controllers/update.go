package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EncargadoEtiquetasRequest struct {
	EncargadoEtiquetas string `json:"encargadoEtiquetasId"`
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

	//filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"encargadoEtiquetasId": responsable.EncargadoEtiquetas,
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
