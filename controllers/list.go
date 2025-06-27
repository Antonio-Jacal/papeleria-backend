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
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterList(c *gin.Context) {

	lista := models.List{}

	if err := c.ShouldBindBodyWithJSON(&lista); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Datos Invalidos"})
		return
	}

	collection := config.GetCollection("pedidos")

	numero, err := utils.GenerateNextNumeroLista(collection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Ocurrio un error en el servidor"})
		return
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
	lista.EncargadoEtiquetas = ""
	lista.StatusForrado = "Por forrar"
	lista.PreparadoPorId = ""

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, lista)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Ocurrio un error, no es posible guardar el documento"})
		return
	} else {
		fmt.Println("Mandamos confirmacion por correo")
		c.JSON(http.StatusOK, gin.H{"Lista confirmada, correo enviado a": lista.Correo})
		return
	}

}

func GetList(c *gin.Context) {
	param := c.Query("grado")
	if param == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'lista' es requerido"})
		return
	}

	filter := bson.M{"grado": param}

	collection := config.GetCollection("listas")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "No existe esa lista"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lista": result})
}
