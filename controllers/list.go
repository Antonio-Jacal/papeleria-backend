package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/Antonio-Jacal/papeleria-backend.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	filterBuilders = map[reflect.Type]func(string, interface{}) bson.D{
		reflect.TypeOf(""): func(key string, value interface{}) bson.D {
			return BuildStringFilter(key, value.(string))
		},
	}
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
		if os.Getenv("ENV") != "production" {
			err := godotenv.Load()
			if err != nil {
				log.Println("No se cargó el archivo .env")
			}
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
	/*
		numeroLista string
		nombreTutor (regex autocomplete) string
		nombreAlumno (regex autocomplete) string
		grado string
		fechaCreacionInicial fecha
		fechaCreacionFinal fecha
		fechaEntregaInicial fecha
		fechaEntregaFinal fecha
		statusLista string
		statusForrado string
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection("pedidos")

	if len(c.Request.URL.Query()) > 0 {
		var mustFilters bson.A

		for key, values := range c.Request.URL.Query() {
			value := values[0]
			if value != "" {
				switch key {
				case "nombreTutor", "nombreAlumno":
					mustFilters = append(mustFilters, BuildAutocompleteFilter(key, value))
				case "numeroLista", "grado", "statusLista", "statusForrado":
					mustFilters = append(mustFilters, BuildStringFilter(key, value))
				}
			}
		}

		var filterRanges bson.A

		if r := buildDateRangeFilter(c.Query("fechaCreacionInicial"), c.Query("fechaCreacionFinal"), "fechaCreacion"); r != nil {
			filterRanges = append(filterRanges, r)
		}
		if r := buildDateRangeFilter(c.Query("fechaEntregaInicial"), c.Query("fechaEntregaFinal"), "fechaEntregaEsperada"); r != nil {
			filterRanges = append(filterRanges, r)
		}

		searchStage := bson.D{}

		if len(filterRanges) > 0 && len(mustFilters) > 0 {
			searchStage = bson.D{{
				Key: "$search", Value: bson.D{
					{Key: "index", Value: "pedidos"},
					{Key: "compound", Value: bson.D{
						{Key: "must", Value: mustFilters},
						{Key: "filter", Value: filterRanges},
					}},
				},
			}}
		} else if len(mustFilters) > 0 {
			searchStage = bson.D{{
				Key: "$search", Value: bson.D{
					{Key: "index", Value: "pedidos"},
					{Key: "compound", Value: bson.D{
						{Key: "must", Value: mustFilters},
					}},
				},
			}}
		} else {
			searchStage = bson.D{{
				Key: "$search", Value: bson.D{
					{Key: "index", Value: "pedidos"},
					{Key: "compound", Value: bson.D{
						{Key: "filter", Value: filterRanges},
					}},
				},
			}}
		}

		pipeline := mongo.Pipeline{
			searchStage,
			bson.D{{Key: "$sort", Value: bson.D{{Key: "fechaCreacion", Value: 1}}}},
		}
		/*
			fmt.Println("Pipeline ejecutado:")
			for _, stage := range pipeline {
				jsonStage, _ := json.MarshalIndent(stage, "", "  ")
				fmt.Println(string(jsonStage))
			}
		*/
		cursor, err := collection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo ejecutar el filtro"})
			return
		}
		defer cursor.Close(ctx)

		var resultados []bson.M
		if err := cursor.All(ctx, &resultados); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron procesar los resultados"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"resultados": resultados})
	} else {
		pipeline := mongo.Pipeline{}
		cursor, err := collection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo ejecutar el filtro"})
			return
		}
		defer cursor.Close(ctx)

		var resultados []bson.M
		if err := cursor.All(ctx, &resultados); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron procesar los resultados"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"resultados": resultados})
	}
}

func BuildStringFilter(key, value string) bson.D {
	return bson.D{
		{Key: "text", Value: bson.D{
			{Key: "query", Value: value},
			{Key: "path", Value: key},
		}},
	}
}

func BuildAutocompleteFilter(key, value string) bson.D {
	return bson.D{{
		Key: "autocomplete", Value: bson.D{
			{Key: "query", Value: value},
			{Key: "path", Value: key},
		},
	}}
}

func buildDateRangeFilter(at, between, campo string) bson.D {
	var t1, t2 *time.Time
	if at != "" {
		t1, _ = utils.ParseTimeParam(at)
	}
	if between != "" {
		t2, _ = utils.ParseTimeParam(between)
	}

	rangeQuery := bson.D{{Key: "path", Value: campo}}

	switch {
	case t1 != nil && t2 != nil:
		if t1.After(*t2) {
			t1, t2 = t2, t1
		}
		rangeQuery = append(rangeQuery, bson.E{Key: "gte", Value: t1})
		rangeQuery = append(rangeQuery, bson.E{Key: "lte", Value: t2})
	case t1 != nil:
		rangeQuery = append(rangeQuery, bson.E{Key: "gte", Value: t1})
	case t2 != nil:
		rangeQuery = append(rangeQuery, bson.E{Key: "lte", Value: t2})
	default:
		return nil
	}

	return bson.D{{Key: "range", Value: rangeQuery}}
}
