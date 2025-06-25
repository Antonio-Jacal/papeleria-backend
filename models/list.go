package models

type list struct {
	ID primitive.ObjectID
	NumeroLista string
	PIN string
	NombreTutor string
	NombreAlumno string
	Correo string
	Telefono string
	Grado string
	FechaCreacion *time.Time
	FechaEntregaEsperada *time.Time
	FechaEntregaReal *time.Time
	EstadoLista string
	Productos map[string]ProductDetail
}

type ProductDetail struct {
	
}

"estadoLista":"Por preparar",
	"productos": {
		"Lápiz": {
			"cantidad":{"$numberInt":"3"},
			"preparado":{"$numberInt":"1"}
		},"Cuaderno":{"cantidad":{"$numberInt":"5"},"preparado":{"$numberInt":"5"}}},"utilesQuitados":{"Mochila":{"$numberInt":"1"},"Juego de Geometría":{"$numberInt":"1"}},"deseaQuitar":true,"faltantes":{"Goma":{"$numberInt":"2"},"Lápiz":{"$numberInt":"1"}},"etiquetasId":"ObjectId","statusEtiquetas":"Sin preparar","encargadoEtiquetasId":"ObjectId","listaForrada":true,"statusForrado":"Por forrar","formaPago":"Efectivo","estaPagado":false,"pagos":[{"monto":{"$numberDouble":"550.5"},"fecha":"ISODate","formaPago":"Efectivo"}],"totalLista":{"$numberDouble":"2000.5"},"totalForrado":{"$numberDouble":"550.5"},"totalGeneral":{"$numberDouble":"2550.5"},"totalPagado":{"$numberDouble":"550.5"},"totalRestante":{"$numberDouble":"2000.0"},"preparadoPorId":"ObjectId","creadoPorId":"ObjectId"}
