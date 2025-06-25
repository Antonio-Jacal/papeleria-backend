package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	ID                   primitive.ObjectID         `json:"id,omitempty" bson:"_id,omitempty"`
	NumeroLista          string                     `json:"numero_lista,omitempty" bson:"numeroLista,omitempty"`
	PIN                  string                     `json:"pin,omitempty" bson:"PIN,omitempty"`
	NombreTutor          string                     `json:"nombre_tutor,omitempty" bson:"nombreTutor,omitempty"`
	NombreAlumno         string                     `json:"nombre_alumno,omitempty" bson:"nombreAlumno,omitempty"`
	Correo               string                     `json:"correo,omitempty" bson:"correo,omitempty"`
	Telefono             string                     `json:"telefono,omitempty" bson:"telefono,omitempty"`
	Grado                string                     `json:"grado,omitempty" bson:"grado,omitempty"`
	FechaCreacion        *time.Time                 `json:"fecha_creacion,omitempty" bson:"fechaCreacion,omitempty"`
	FechaEntregaEsperada *time.Time                 `json:"fecha_entrega_esperada,omitempty" bson:"fechaEntregaEsperada,omitempty"`
	FechaEntregaReal     *time.Time                 `json:"fecha_entrega_real,omitempty" bson:"fechaEntregaReal,omitempty"`
	EstadoLista          string                     `json:"estado_lista,omitempty" bson:"estadoLista,omitempty"`
	Productos            map[string]ProductoDetalle `json:"productos,omitempty" bson:"productos,omitempty"`
	UtilesQuitados       map[string]int             `json:"utiles_quitados,omitempty" bson:"utilesQuitados,omitempty"`
	DeseaQuitar          bool                       `json:"desea_quitar,omitempty" bson:"deseaQuitar,omitempty"`
	Faltantes            map[string]int             `json:"faltantes,omitempty" bson:"faltantes,omitempty"`
	EtiquetasId          primitive.ObjectID         `json:"etiquetas_id,omitempty" bson:"etiquetasId,omitempty"`
	StatusEtiquetas      string                     `json:"status_etiquetas,omitempty" bson:"statusEtiquetas,omitempty"`
	EncargadoEtiquetasID primitive.ObjectID         `json:"encargado_etiquetas_id,omitempty" bson:"encargadoEtiquetasId,omitempty"`
	ListaForrada         bool                       `json:"lista_forrada,omitempty" bson:"listaForrada,omitempty"`
	StatusForrado        string                     `json:"status_forrado,omitempty" bson:"statusForrado,omitempty"`
	FormaPago            string                     `json:"forma_pago,omitempty" bson:"formaPago,omitempty"`
	EstaPagado           bool                       `json:"esta_pagado,omitempty" bson:"estaPagado,omitempty"`
	Pagos                []Pago                     `json:"pagos,omitempty" bson:"pagos,omitempty"`
	TotalLista           float32                    `json:"total_lista,omitempty" bson:"totalLista,omitempty"`
	TotalForrado         float32                    `json:"total_forrado,omitempty" bson:"totalForrado,omitempty"`
	TotalGeneral         float32                    `json:"total_general,omitempty" bson:"totalGeneral,omitempty"`
	TotalPagado          float32                    `json:"total_pagado,omitempty" bson:"totalPagado,omitempty"`
	TotalRestante        float32                    `json:"total_restante,omitempty" bson:"totalRestante,omitempty"`
	PreparadoPorId       primitive.ObjectID         `json:"preparado_por_id,omitempty" bson:"preparadoPorId,omitempty"`
	CreadoPorId          primitive.ObjectID         `json:"creado_por_id,omitempty" bson:"creadoPorId,omitempty"`
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
