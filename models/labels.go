package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LabelResumen struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	NumeroLista          string             `bson:"numeroLista"`
	NombreAlumno         string             `bson:"nombreAlumno"`
	Grado                string             `bson:"grado"`
	FechaEntregaEsperada *time.Time         `bson:"fechaEntregaEsperada"`
	EtiquetasPersonaje   string             `bson:"etiquetasPersonaje"`
	StatusEtiquetas      string             `bson:"statusEtiquetas"`
	EtiquetasGrandes     bool               `bson:"etiquetasGrandes"`
	EtiquetasMedianas    bool               `bson:"etiquetasMedianas"`
	EtiquetasChicas      bool               `bson:"etiquetasChicas"`
	EncargadoEtiquetasId string             `bson:"encargadoEtiquetasId"`
	StatusForrado        string             `bson:"statusForrado"`
}
