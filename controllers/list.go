package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/Antonio-Jacal/papeleria-backend.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterList(c *gin.Context) {
	var lista models.List

	if err := c.ShouldBindJSON(&lista); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   "Datos inválidos",
			"Detalle": err.Error(),
		})
		return
	}

	// Validación mejorada
	missingFields := []string{}
	if lista.NombreTutor == "" {
		missingFields = append(missingFields, "nombreTutor")
	}
	if lista.NombreAlumno == "" {
		missingFields = append(missingFields, "nombreAlumno")
	}
	if lista.Correo == "" {
		missingFields = append(missingFields, "correo")
	}
	if lista.Grado == "" {
		missingFields = append(missingFields, "grado")
	}
	if lista.Telefono == "" {
		missingFields = append(missingFields, "numero")
	}

	if len(missingFields) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error":  "Campos obligatorios faltantes",
			"Campos": missingFields,
		})
		return
	}
	// Asegura los valores por defecto
	if lista.FechaCreacion == nil {
		loc, _ := time.LoadLocation("America/Mexico_City")
		now := time.Now().In(loc)
		lista.FechaCreacion = &now
	}

	// Configuración de campos automáticos
	lista.EstadoLista = "Por preparar"
	lista.PIN = utils.GeneratePin()

	// Configuración condicional
	if lista.ListaForrada {
		lista.StatusForrado = "Por forrar"
		lista.StatusEtiquetas = "Por hacer"
		lista.EtiquetasChicas = true
		lista.EtiquetasGrandes = true
		lista.EtiquetasMedianas = true
	} else {
		lista.StatusForrado = "No aplica"
		lista.StatusEtiquetas = "No aplica"
		lista.EtiquetasChicas = false
		lista.EtiquetasGrandes = false
		lista.EtiquetasMedianas = false
	}

	// Generar número de lista
	collection := config.GetCollection("pedidos")
	numero, err := utils.GenerateNextNumeroLista(collection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error generando número de lista"})
		return
	}
	lista.NumeroLista = numero

	// Insertar en MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//fmt.Println(lista)

	_, err = collection.InsertOne(ctx, lista)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Ocurrio un error, no es posible guardar el documento"})
		return
	} else {
		fmt.Println("Mandamos confirmacion por correo")
		err := godotenv.Load()
		if err != nil {
			log.Println("No se pudo cargar el .env, usando variables del sistema")
		}
		to := []string{lista.Correo}
		subject := fmt.Sprintf("Confirmación de pedido %s", lista.NumeroLista)
		html := fmt.Sprintf(`
<html>
  <body style="font-family: sans-serif; color: #333;">
    <div style="max-width: 600px; margin: auto; border: 1px solid #ddd; padding: 30px; border-radius: 10px;">
      <h2 style="color: #2196F3;">Confirmación de Pedido: %s</h2>
      <p>¡Hola <strong>%s</strong>!</p>
      <p>El pedido para <strong>%s</strong> (Grado: <strong>%s</strong>) ha sido registrado exitosamente.</p>

      <p><strong>Detalles del pedido:</strong></p>
      <ul>
        <li><strong>Número de lista:</strong> %s</li>
        <li><strong>Fecha de creación:</strong> %s</li>
        <li><strong>Fecha estimada de entrega:</strong> %s</li>
        <li><strong>Etiquetas:</strong> %s</li>
      </ul>

      <p><strong>Productos solicitados:</strong></p>
      %s

      <p><strong>Útiles quitados:</strong></p>
      %s

      <hr style="margin: 20px 0;" />

      <p><strong>Total a pagar:</strong> $%.2f MXN</p>
      <p><strong>Total pagado:</strong> $%.2f MXN</p>
      <p><strong>Total restante:</strong> $%.2f MXN</p>

      <p style="font-size: 14px; color: #888; margin-top: 30px;">
        Gracias por confiar en nosotros.<br>
        <em>Equipo de Papelería Nina's</em>
      </p>
    </div>
  </body>
</html>
`,
			lista.NumeroLista,
			lista.NombreTutor,
			lista.NombreAlumno,
			lista.Grado,
			lista.NumeroLista,
			utils.FormatDate(lista.FechaCreacion),
			utils.FormatDate(lista.FechaEntregaEsperada),
			lista.EtiquetasPersonaje,
			utils.BuildProductosHTML(lista.Productos),
			utils.BuildUtilesQuitadosHTML(lista.UtilesQuitados),
			lista.TotalGeneral,
			lista.TotalPagado,
			lista.TotalRestante,
		)

		err = utils.SendHTMLEmail(to, subject, html)
		if err != nil {
			log.Fatal("Fallo al enviar correo:", err)
		}

		log.Println("Correo enviado exitosamente.")
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

func GetListWithFilters(c *gin.Context) {

}
