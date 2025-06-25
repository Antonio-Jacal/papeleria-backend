package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/routes"
)

func main() {
	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error cargando .env")
	}

	// Conecta a MongoDB
	config.ConnectDB()

	// Inicializa servidor
	r := gin.Default()
	routes.SetupRoutes(r)

	log.Println("Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
