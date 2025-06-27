package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/Antonio-Jacal/papeleria-backend.git/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

	lista.FechaEntregaReal = nil
	loc, _ := time.LoadLocation("America/Mexico_City")
	now := time.Now().In(loc)
	lista.FechaCreacion = &now
	lista.EstadoLista = "Por preparar" // Por preparar, prerada, lista, Con Faltantes
	lista.Faltantes = nil
	if lista.EtiquetasPersonaje == "" {
		lista.EtiquetasPersonaje = "Por confirmar"
	}
	if lista.ListaForrada {
		lista.StatusEtiquetas = "Por hacer"
		lista.EtiquetasGrandes = true
		lista.EtiquetasMedinas = true
		lista.EtiquetasChicas = true
	} else {
		lista.StatusEtiquetas = "No aplica"
		lista.EtiquetasGrandes = false
		lista.EtiquetasMedinas = false
		lista.EtiquetasChicas = false
	}
	lista.StatusForrado = "Por forrar"
	lista.PreparadoPorId = ""

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, lista)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Ocurrio un error, no es posible guardar el documento"})
	} else {
		fmt.Println("Mandamos confirmacion por correo")
		c.JSON(http.StatusOK, gin.H{"Lista confirmada, correo enviado a": lista.Correo})
	}

}

func GetList(c *gin.Context) {
	param := c.Query("lista")
	filter := bson.M{}
	filter["grado"] = param

	collection := config.GetCollection("listas")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println(filter)
	lista := collection.FindOne(ctx, bson.M{"grado": "Primaria 2"})
	if lista == nil {
		c.JSON(http.StatusOK, "No existe esa lista")
	}
	c.JSON(http.StatusOK, gin.H{"lista": lista})

}
