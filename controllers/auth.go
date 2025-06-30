package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Antonio-Jacal/papeleria-backend.git/config"
	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/Antonio-Jacal/papeleria-backend.git/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(c *gin.Context) {
	var input struct {
		Nombre   string `json:"nombre"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Rol      string `json:"rol"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Validar rol permitido
	if input.Rol != "worker" && input.Rol != "admin" && input.Rol != "develop" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rol inválido"})
		return
	}

	// Verifica si ya existe
	collection := config.GetCollection("usuarios")
	var existing models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El usuario ya existe"})
		return
	}

	hashedPwd, _ := utils.HashPassword(input.Password)

	user := models.User{
		ID:       primitive.NewObjectID(),
		Nombre:   input.Nombre,
		Email:    input.Email,
		Password: hashedPwd,
		Rol:      input.Rol,
	}

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}

	// Enviar correo de bienvenida
	htmlBody := fmt.Sprintf(`
<html>
  <body style="font-family: sans-serif; color: #333;">
    <div style="max-width: 600px; margin: auto; border: 1px solid #ddd; padding: 30px; border-radius: 10px;">
      <h2 style="color: #4CAF50;">¡Bienvenido al equipo, %s! 🎉</h2>
      <p style="font-size: 16px;">
        ¿Estás listo para una nueva temporada de <strong>papelería</strong>?<br><br>
        Estamos muy emocionados de tenerte con nosotros.
      </p>
      <p style="font-size: 16px; margin-top: 30px;">
        <strong>Tu contraseña para acceder a la plataforma es:</strong>
      </p>
      <div style="background-color: #f2f2f2; padding: 15px; border-radius: 8px; font-size: 18px; font-weight: bold; text-align: center;">
        %s
      </div>
      <p style="font-size: 14px;">¡Mucho éxito!<br>El equipo de Papelería</p>
    </div>
  </body>
</html>
`, input.Nombre, input.Password)

	go func() {
		err := utils.SendHTMLEmail([]string{input.Email}, "🎒 Bienvenido a la plataforma de papelería Nina's", htmlBody)
		if err != nil {
			log.Printf("Error al enviar correo de bienvenida a %s: %v", input.Email, err)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{"mensaje": "Usuario registrado correctamente"})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var user models.User
	collection := config.GetCollection("usuarios")
	err := collection.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil || !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), user.Rol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  token,
		"userId": user.ID.Hex(),
		"rol":    user.Rol,
		"name":   user.Nombre,
	})
}
