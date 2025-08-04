package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	ID                   primitive.ObjectID         `json:"id,omitempty" bson:"_id,omitempty"`
	NumeroLista          string                     `json:"numeroLista,omitempty" bson:"numeroLista,omitempty"`
	PIN                  string                     `json:"pin,omitempty" bson:"pin,omitempty"` // en bson estaba "PIN" (mayúsculas), mejor consistente con JSON
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
	DeseaQuitar          bool                       `json:"deseaQuitar,omitempty" bson:"deseaQuitar"`
	Faltantes            map[string]int             `json:"faltantes,omitempty" bson:"faltantes,omitempty"`
	ListaForrada         bool                       `json:"listaForrada,omitempty" bson:"listaForrada"`
	EtiquetasPersonaje   string                     `json:"etiquetasPersonaje,omitempty" bson:"etiquetasPersonaje,omitempty"`
	StatusEtiquetas      string                     `json:"statusEtiquetas,omitempty" bson:"statusEtiquetas,omitempty"` // corregí tag (mayúscula y typo "Somitempty")
	EtiquetasGrandes     bool                       `json:"etiquetasGrandes,omitempty" bson:"etiquetasGrandes"`
	EtiquetasMedianas    bool                       `json:"etiquetasMedianas,omitempty" bson:"etiquetasMedianas"`
	EtiquetasChicas      bool                       `json:"etiquetasChicas,omitempty" bson:"etiquetasChicas"`
	EncargadoEtiquetas   string                     `json:"encargadoEtiquetasId,omitempty" bson:"encargadoEtiquetasId,omitempty"` // ajustado para que coincida
	StatusForrado        string                     `json:"statusForrado,omitempty" bson:"statusForrado,omitempty"`
	StatusForradoLibros  string                     `json:"statusForradoLibros,omitempty" bson:"statusForradoLibros,omitempty"`
	StatusForradoUtiles  string                     `json:"statusForradoUtiles,omitempty" bosn:"statusForradoUtiles,omitempty"`
	FormaPago            string                     `json:"formaPago,omitempty" bson:"formaPago,omitempty"`
	EstaPagado           bool                       `json:"estaPagado,omitempty" bson:"estaPagado"`
	Pagos                []Pago                     `json:"pagos,omitempty" bson:"pagos,omitempty"`
	TotalLista           float64                    `json:"totalLista,omitempty" bson:"totalLista,omitempty"`
	TotalForrado         float64                    `json:"totalForrado,omitempty" bson:"totalForrado,omitempty"`
	TotalGeneral         float64                    `json:"totalGeneral,omitempty" bson:"totalGeneral,omitempty"`
	TotalPagado          float64                    `json:"totalPagado,omitempty" bson:"totalPagado,omitempty"`
	TotalRestante        float64                    `json:"totalRestante,omitempty" bson:"totalRestante,omitempty"`
	PreparadoPorId       string                     `json:"preparadoPorId,omitempty" bson:"preparadoPorId,omitempty"`
	Comentarios          string                     `json:"comentarios,omitempty" bson:"comentarios,omitempty"`
}

type ProductoDetalle struct {
	Cantidad  int `bson:"cantidad,omitempty" json:"cantidad,omitempty"`
	Preparado int `bson:"preparado,omitempty" json:"preparado,omitempty"`
}

type Pago struct {
	Monto     float64    `bson:"monto,omitempty" json:"monto,omitempty"`
	Fecha     *time.Time `bson:"fecha,omitempty" json:"fecha,omitempty"`
	FormaPago string     `bson:"formaPago,omitempty" json:"formaPago,omitempty"` // corregí json tag para coincidir con bson
}

type FilterList struct {
	//Filtros:
	//Numero de lista
	//nombre Tutor (regex)
	//nombre Alumno (regex)
	//grado
	//fecha Creacion
	//Fecha de entrega
	//Status lista
	//Status forrado
	NumLista      string     `json:"num_lista,omitempty"`
	NombreTutor   string     `json:"nombre_tutor,omitempty"`
	NombreAlumno  string     `json:"nombre_alumno,omitempty"`
	Grado         string     `json:"grado,omitempty"`
	StatusLista   string     `json:"status_lista,omitempty"`
	StatusForrado string     `json:"status_forrado,omitempty"`
	FechaCreacion *time.Time `json:"fecha_creacion,omitempty"`
	FechaEntrega  *time.Time `json:"fecha_entrega,omitempty"`
}
