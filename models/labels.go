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
	Comentarios          string             `bson:"comentarios"`
	FechaEntregaEsperada *time.Time         `bson:"fechaEntregaEsperada"`
	EtiquetasPersonaje   string             `bson:"etiquetasPersonaje"`
	StatusEtiquetas      string             `bson:"statusEtiquetas"`
	EtiquetasGrandes     bool               `bson:"etiquetasGrandes"`
	EtiquetasMedianas    bool               `bson:"etiquetasMedianas"`
	EtiquetasChicas      bool               `bson:"etiquetasChicas"`
	NumPaquete           int                `bson:"numPaquete"`
	Tipografia           int                `bson:"tipografia"`
	Marco                string             `bson:"marco"`
	Patron               int                `bson:"patron"`
	Acomodo              int                `bson:"acomodo"`
	EncargadoEtiquetasId string             `bson:"encargadoEtiquetasId"`
	StatusForrado        string             `bson:"statusForrado"`
}
