package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	ID                   primitive.ObjectID         `json:"id,omitempty" bson:"_id,omitempty"`
	NumeroLista          string                     `json:"numeroLista,omitempty" bson:"numeroLista,omitempty"`
	PIN                  string                     `json:"pin,omitempty" bson:"PIN,omitempty"`
	NombreTutor          string                     `json:"nombreTutor,omitempty" bson:"nombreTutor,omitempty"`
	NombreAlumno         string                     `json:"nombreAlumno,omitempty" bson:"nombreAlumno,omitempty"`
	Correo               string                     `json:"correo,omitempty" bson:"correo,omitempty"`
	Telefono             string                     `json:"telefono,omitempty" bson:"telefono,omitempty"`
	Grado                string                     `json:"grado,omitempty" bson:"grado,omitempty"`
	FechaCreacion        *time.Time                 `json:"fechaCreacion,omitempty" bson:"fechaCreacion,omitempty"`
	FechaEntregaEsperada *time.Time                 `json:"fechaEntregaEsperada,omitempty" bson:"fechaEntregaEsperada,omitempty"`
	FechaEntregaReal     *time.Time                 `json:"fechaEntregaReal,omitempty" bson:"fechaEntregaReal,omitempty"`
	EstadoLista          string                     `json:"estadoLista,omitempty" bson:"estadoLista,omitempty"`
	Productos            map[string]ProductoDetalle `json:"productos,omitempty" bson:"productos,omitempty"`
	UtilesQuitados       map[string]int             `json:"utilesQuitados,omitempty" bson:"utilesQuitados,omitempty"`
	DeseaQuitar          bool                       `json:"deseaQuitar" bson:"deseaQuitar"`
	Faltantes            map[string]int             `json:"faltantes,omitempty" bson:"faltantes,omitempty"`
	ListaForrada         bool                       `json:"listaForrada" bson:"listaForrada"`
	EtiquetasPersonaje   string                     `json:"etiquetasPersonaje,omitempty" bson:"etiquetasPersonaje,omitempty"`
	StatusEtiquetas      string                     `json:"status_etiquetas,omitempty" bson:"statusEtiquetas,Somitempty"`
	EtiquetasGrandes     bool                       `json:"etiquetasGrandes" bson:"etiquetasGrandes"`
	EtiquetasMedianas    bool                       `json:"etiquetasMedianas" bson:"etiquetasMedianas"`
	EtiquetasChicas      bool                       `json:"etiquetasChicas" bson:"etiquetasChicas"`
	EncargadoEtiquetas   string                     `json:"encargadoEtiquetas_id,omitempty" bson:"encargadoEtiquetasId,omitempty"`
	StatusForrado        string                     `json:"statusForrado,omitempty" bson:"statusForrado,omitempty"`
	FormaPago            string                     `json:"formaPago,omitempty" bson:"formaPago,omitempty"`
	EstaPagado           bool                       `json:"estaPagado" bson:"estaPagado"`
	Pagos                []Pago                     `json:"pagos,omitempty" bson:"pagos,omitempty"`
	TotalLista           float64                    `json:"totalLista,omitempty" bson:"totalLista,omitempty"`
	TotalForrado         float64                    `json:"totalForrado,omitempty" bson:"totalForrado,omitempty"`
	TotalGeneral         float64                    `json:"totalGeneral,omitempty" bson:"totalGeneral,omitempty"`
	TotalPagado          float64                    `json:"totalPagado,omitempty" bson:"totalPagado,omitempty"`
	TotalRestante        float64                    `json:"totalRestante,omitempty" bson:"totalRestante,omitempty"`
	PreparadoPorId       string                     `json:"preparadoPorId,omitempty" bson:"preparadoPorId,omitempty"`
}

type ProductoDetalle struct {
	Cantidad  int `bson:"cantidad" json:"cantidad,omitempty"`
	Preparado int `bson:"preparado" json:"preparado,omitempty"`
}

type Pago struct {
	Monto     float64            `bson:"monto" json:"monto,omitempty"`
	Fecha     primitive.DateTime `bson:"fecha" json:"fecha,omitempty"`
	FormaPago string             `bson:"formaPago" json:"forma_pago,omitempty"`
}
