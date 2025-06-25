package controllers

import (
	"net/http"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/Antonio-Jacal/papeleria-backend.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterList(c *gin.Context) {

	lista := models.List{}

	if err := c.ShouldBindBodyWithJSON(&lista); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos Invalidos"})
	}

	collection := config.GetCollection("pedidos")

	numero, err := utils.GenerateNextNumeroLista(collection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Ocurrio un error en el servidor"})
	}

	lista.NumeroLista = numero
	lista.PIN = utils.GeneratePin()

}
