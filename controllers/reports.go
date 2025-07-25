package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	result, err := collection.Find(ctx, bson.M{"listaForrada": true}, options.Find().SetProjection(getProjection()).SetSort(bson.D{{Key: "numeroLista", Value: 1}}))

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

func GetUrgent(c *gin.Context) {
	filter := bson.M{
		"estadoLista": bson.M{"$ne": "Entregada"},
		"fechaEntregaEsperada": bson.M{
			"$gte": time.Now(),
			"$lte": time.Now().AddDate(0, 0, 3),
		},
	}
	collection := config.GetCollection("pedidos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(setProjectionUrgent()).SetSort(bson.M{"fechaEntregaEsperada": 1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": err})
		return
	}

	type resultUrgent struct {
		ID                   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
		NumeroLista          string             `json:"numeroLista,omitempty"`
		NombreAlumno         string             `json:"nombreAlumno,omitempty"`
		NombreTutor          string             `json:"nombreTutor,omitempty"`
		Grado                string             `json:"grado,omitempty"`
		EstadoLista          string             `json:"estadoLista,omitempty"`
		StatusForrado        string             `json:"statusForrado,omitempty"`
		StatusEtiquetas      string             `json:"statusEtiquetas,omitempty"`
		FechaEntregaEsperada *time.Time         `json:"fechaEntregaEsperada,omitempty"`
	}

	var result []resultUrgent
	for cursor.Next(ctx) {
		var item resultUrgent
		if err := cursor.Decode(&item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error al decodificar documento:": err})
			return
		}
		result = append(result, item)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error en el cursor:": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resutados": result})
}

func setProjectionUrgent() bson.M {
	return bson.M{
		"_id":                  1,
		"numeroLista":          1,
		"nombreAlumno":         1,
		"nombreTutor":          1,
		"grado":                1,
		"estadoLista":          1,
		"statusForrado":        1,
		"statusEtiquetas":      1,
		"fechaEntregaEsperada": 1,
	}
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

func GetSummaryFaltantes(c *gin.Context) {
	collection := config.GetCollection("pedidos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"estadoLista": "Con faltantes"}, options.Find().SetProjection(bson.M{"faltantes": 1, "grado": 1, "numeroLista": 1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer cursor.Close(ctx)

	type ListFaltantes struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		NumeroLista string             `bson:"numeroLista,omitempty" json:"numeroLista,omitempty"`
		Faltantes   map[string]int     `bson:"faltantes,omitempty" json:"faltantes,omitempty"`
		Grado       string             `bson:"grado,omitempty" json:"grado,omitempty"`
	}
	var results []ListFaltantes

	if err := cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	faltantesPorGrado := map[string]map[string]int{}

	for _, list := range results {
		switch list.Grado {
		case "Peques":
			if faltantesPorGrado["Peques"] == nil {
				faltantesPorGrado["Peques"] = make(map[string]int)
				faltantesPorGrado["Peques"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Peques"][item] += cantidad
			}
			faltantesPorGrado["Peques"]["TotalListasFaltantes"] += 1
		case "Preescolar 1":
			if faltantesPorGrado["Preescolar 1"] == nil {
				faltantesPorGrado["Preescolar 1"] = make(map[string]int)
				faltantesPorGrado["Preescolar 1"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Preescolar 1"][item] += cantidad
			}
			faltantesPorGrado["Preescolar 1"]["TotalListasFaltantes"] += 1
		case "Preescolar 2":
			if faltantesPorGrado["Preescolar 2"] == nil {
				faltantesPorGrado["Preescolar 2"] = make(map[string]int)
				faltantesPorGrado["Preescolar 2"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Preescolar 2"][item] += cantidad
			}
			faltantesPorGrado["Preescolar 2"]["TotalListasFaltantes"] += 1
		case "Preescolar 3":
			if faltantesPorGrado["Preescolar 3"] == nil {
				faltantesPorGrado["Preescolar 3"] = make(map[string]int)
				faltantesPorGrado["Preescolar 3"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Preescolar 3"][item] += cantidad
			}
			faltantesPorGrado["Preescolar 3"]["TotalListasFaltantes"] += 1
		case "Primaria 1":
			if faltantesPorGrado["Primaria 1"] == nil {
				faltantesPorGrado["Primaria 1"] = make(map[string]int)
				faltantesPorGrado["Primaria 1"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Primaria 1"][item] += cantidad
			}
			faltantesPorGrado["Primaria 1"]["TotalListasFaltantes"] += 1
		case "Primaria 2":
			if faltantesPorGrado["Primaria 2"] == nil {
				faltantesPorGrado["Primaria 2"] = make(map[string]int)
				faltantesPorGrado["Primaria 2"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Primaria 2"][item] += cantidad
			}
			faltantesPorGrado["Primaria 2"]["TotalListasFaltantes"] += 1
		case "Primaria 3":
			if faltantesPorGrado["Primaria 3"] == nil {
				faltantesPorGrado["Primaria 3"] = make(map[string]int)
				faltantesPorGrado["Primaria 3"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Primaria 3"][item] += cantidad
			}
			faltantesPorGrado["Primaria 3"]["TotalListasFaltantes"] += 1
		case "Primaria 4":
			if faltantesPorGrado["Primaria 4"] == nil {
				faltantesPorGrado["Primaria 4"] = make(map[string]int)
				faltantesPorGrado["Primaria 4"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Primaria 4"][item] += cantidad
			}
			faltantesPorGrado["Primaria 4"]["TotalListasFaltantes"] += 1
		case "Primaria 5":
			if faltantesPorGrado["Primaria 5"] == nil {
				faltantesPorGrado["Primaria 5"] = make(map[string]int)
				faltantesPorGrado["Primaria 5"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Primaria 5"][item] += cantidad
			}
			faltantesPorGrado["Primaria 5"]["TotalListasFaltantes"] += 1
		case "Primaria 6":
			if faltantesPorGrado["Primaria 6"] == nil {
				faltantesPorGrado["Primaria 6"] = make(map[string]int)
				faltantesPorGrado["Primaria 6"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Primaria 6"][item] += cantidad
			}
			faltantesPorGrado["Primaria 6"]["TotalListasFaltantes"] += 1
		case "Secundaria 1":
			if faltantesPorGrado["Secundaria 1"] == nil {
				faltantesPorGrado["Secundaria 1"] = make(map[string]int)
				faltantesPorGrado["Secundaria 1"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Secundaria 1"][item] += cantidad
			}
			faltantesPorGrado["Secundaria 1"]["TotalListasFaltantes"] += 1
		case "Secundaria 2":
			if faltantesPorGrado["Secundaria 2"] == nil {
				faltantesPorGrado["Secundaria 2"] = make(map[string]int)
				faltantesPorGrado["Secundaria 2"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Secundaria 2"][item] += cantidad
			}
			faltantesPorGrado["Secundaria 2"]["TotalListasFaltantes"] += 1
		case "Secundaria 3":
			if faltantesPorGrado["Secundaria 3"] == nil {
				faltantesPorGrado["Secundaria 3"] = make(map[string]int)
				faltantesPorGrado["Secundaria 3"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Secundaria 3"][item] += cantidad
			}
			faltantesPorGrado["Secundaria 3"]["TotalListasFaltantes"] += 1
		default:
			if faltantesPorGrado["Desconocido"] == nil {
				faltantesPorGrado["Desconocido"] = make(map[string]int)
				faltantesPorGrado["Desconocido"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Faltantes {
				faltantesPorGrado["Desconocido"][item] += cantidad
			}
			faltantesPorGrado["Desconocido"]["TotalListasFaltantes"] += 1
		}

	}

	cursor, err = collection.Find(ctx, bson.M{"estadoLista": "Por preparar"}, options.Find().SetProjection(bson.M{"productos": 1, "grado": 1, "numeroLista": 1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer cursor.Close(ctx)

	type ProductoDetalle struct {
		Cantidad  int `bson:"cantidad,omitempty" json:"cantidad,omitempty"`
		Preparado int `bson:"preparado,omitempty" json:"preparado,omitempty"`
	}

	type ListPorPreparar struct {
		ID          primitive.ObjectID         `bson:"_id,omitempty" json:"id,omitempty"`
		NumeroLista string                     `bson:"numeroLista,omitempty" json:"numeroLista,omitempty"`
		Productos   map[string]ProductoDetalle `bson:"productos,omitempty" json:"productos,omitempty"`
		Grado       string                     `bson:"grado,omitempty" json:"grado,omitempty"`
	}

	var resultsPorPreparar []ListPorPreparar

	if err := cursor.All(ctx, &resultsPorPreparar); err != nil {
		log.Fatal(err)
	}

	for _, list := range resultsPorPreparar {
		switch list.Grado {
		case "Peques":
			if faltantesPorGrado["Peques"] == nil {
				faltantesPorGrado["Peques"] = make(map[string]int)
				faltantesPorGrado["Peques"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Peques"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Peques"]["TotalListasCompletas"] += 1
		case "Preescolar 1":
			if faltantesPorGrado["Preescolar 1"] == nil {
				faltantesPorGrado["Preescolar 1"] = make(map[string]int)
				faltantesPorGrado["Preescolar 1"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Preescolar 1"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Preescolar 1"]["TotalListasCompletas"] += 1
		case "Preescolar 2":
			if faltantesPorGrado["Preescolar 2"] == nil {
				faltantesPorGrado["Preescolar 2"] = make(map[string]int)
				faltantesPorGrado["Preescolar 2"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Preescolar 2"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Preescolar 2"]["TotalListasCompletas"] += 1
		case "Preescolar 3":
			if faltantesPorGrado["Preescolar 3"] == nil {
				faltantesPorGrado["Preescolar 3"] = make(map[string]int)
				faltantesPorGrado["Preescolar 3"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Preescolar 3"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Preescolar 3"]["TotalListasCompletas"] += 1
		case "Primaria 1":
			if faltantesPorGrado["Primaria 1"] == nil {
				faltantesPorGrado["Primaria 1"] = make(map[string]int)
				faltantesPorGrado["Primaria 1"]["TotalListasFaltantes"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Primaria 1"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Primaria 1"]["TotalListasCompletas"] += 1
		case "Primaria 2":
			if faltantesPorGrado["Primaria 2"] == nil {
				faltantesPorGrado["Primaria 2"] = make(map[string]int)
				faltantesPorGrado["Primaria 2"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Primaria 2"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Primaria 2"]["TotalListasCompletas"] += 1
		case "Primaria 3":
			if faltantesPorGrado["Primaria 3"] == nil {
				faltantesPorGrado["Primaria 3"] = make(map[string]int)
				faltantesPorGrado["Primaria 3"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Primaria 3"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Primaria 3"]["TotalListasCompletas"] += 1
		case "Primaria 4":
			if faltantesPorGrado["Primaria 4"] == nil {
				faltantesPorGrado["Primaria 4"] = make(map[string]int)
				faltantesPorGrado["Primaria 4"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Primaria 4"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Primaria 4"]["TotalListasCompletas"] += 1
		case "Primaria 5":
			if faltantesPorGrado["Primaria 5"] == nil {
				faltantesPorGrado["Primaria 5"] = make(map[string]int)
				faltantesPorGrado["Primaria 5"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Primaria 5"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Primaria 5"]["TotalListasCompletas"] += 1
		case "Primaria 6":
			if faltantesPorGrado["Primaria 6"] == nil {
				faltantesPorGrado["Primaria 6"] = make(map[string]int)
				faltantesPorGrado["Primaria 6"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Primaria 6"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Primaria 6"]["TotalListasCompletas"] += 1
		case "Secundaria 1":
			if faltantesPorGrado["Secundaria 1"] == nil {
				faltantesPorGrado["Secundaria 1"] = make(map[string]int)
				faltantesPorGrado["Secundaria 1"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Secundaria 1"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Secundaria 1"]["TotalListasCompletas"] += 1
		case "Secundaria 2":
			if faltantesPorGrado["Secundaria 2"] == nil {
				faltantesPorGrado["Secundaria 2"] = make(map[string]int)
				faltantesPorGrado["Secundaria 2"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Secundaria 2"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Secundaria 2"]["TotalListasCompletas"] += 1
		case "Secundaria 3":
			if faltantesPorGrado["Secundaria 3"] == nil {
				faltantesPorGrado["Secundaria 3"] = make(map[string]int)
				faltantesPorGrado["Secundaria 3"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Secundaria 3"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Secundaria 3"]["TotalListasCompletas"] += 1
		default:
			if faltantesPorGrado["Desconocido"] == nil {
				faltantesPorGrado["Desconocido"] = make(map[string]int)
				faltantesPorGrado["Desconocido"]["TotalListasCompletas"] = 0
			}
			for item, cantidad := range list.Productos {
				faltantesPorGrado["Desconocido"][item] += cantidad.Cantidad
			}
			faltantesPorGrado["Desconocido"]["TotalListasCompletas"] += 1
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Resumen de faltantes generado correctamente",
		"data": gin.H{
			"faltantes": faltantesPorGrado,
		},
	})

}
