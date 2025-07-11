package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Antonio-Jacal/papeleria-backend.git/models"
)

func SendMessageFromWhatsapp(recipientPhone string) (bool, error) {
	ACCESS_TOKEN := os.Getenv("TOKEN_ACCESS_WHATSAPP")
	PHONENUMBER_ID := os.Getenv("PHONE_NUMBER_ID")

	recipientPhone = fmt.Sprintf("521%s", recipientPhone)

	msg := models.MessageRequest{
		MessagingProduct: "whatsapp",
		To:               recipientPhone,
		Type:             "text",
		Text:             models.MessageText{Body: "¡Hola! Este mensaje fue enviado desde Go con la API de WhatsApp Cloud de Meta 🚀"},
	}

	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", PHONENUMBER_ID)

	body, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error al convertir a JSON:")
		return false, err
	}

	// Construye y envía la solicitud
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error al crear la solicitud:")
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ACCESS_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al enviar la solicitud:")
		return false, err
	}
	defer resp.Body.Close()

	// Muestra el resultado
	fmt.Println("Código de respuesta:", resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("✅ Mensaje enviado correctamente.")
		return true, nil
	} else {
		fmt.Println("⚠️ Error al enviar el mensaje.")
		return false, nil
	}
}
