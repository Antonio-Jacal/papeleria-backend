package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/joho/godotenv"
)

func SendMessageFromWhatsapp(datosLista models.List) (bool, error) {

	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No se cargó el archivo .env")
		}
	}

	ACCESS_TOKEN := os.Getenv("TOKEN_ACCESS_WHATSAPP")
	PHONENUMBER_ID := os.Getenv("PHONE_NUMBER_ID")
	URL_ARCHIVO := "https://rgajiduoagnlivrxfthm.supabase.co/storage/v1/object/public/pedidos//lista_ejemplo.pdf"

	recipientPhone := fmt.Sprintf("52%s", strings.ReplaceAll(datosLista.Telefono, " ", ""))

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                recipientPhone,
		"type":              "template",
		"template": map[string]interface{}{
			"name": "purchase_receipt_1",
			"language": map[string]string{
				"code": "es",
			},
			"components": []interface{}{
				map[string]interface{}{
					"type": "header",
					"parameters": []interface{}{
						map[string]interface{}{
							"type": "document",
							"document": map[string]string{
								"link":     URL_ARCHIVO,
								"filename": "ListaEscolar.pdf",
							},
						},
					},
				},
				map[string]interface{}{ // body con 8 parámetros de texto
					"type": "body",
					"parameters": []interface{}{
						map[string]interface{}{"type": "text", "text": datosLista.NombreTutor},
						map[string]interface{}{"type": "text", "text": datosLista.NumeroLista},
						map[string]interface{}{"type": "text", "text": datosLista.NombreAlumno},
						map[string]interface{}{"type": "text", "text": datosLista.Grado},
						map[string]interface{}{"type": "text", "text": totalFormat(datosLista.TotalGeneral)},
						map[string]interface{}{"type": "text", "text": totalFormat(datosLista.TotalPagado)},
						map[string]interface{}{"type": "text", "text": totalFormat(datosLista.TotalRestante)},
						map[string]interface{}{"type": "text", "text": dateFormat(datosLista.FechaEntregaEsperada.String()[:10])},
						map[string]interface{}{"type": "text", "text": datosLista.Correo},
					},
				},
			},
		},
	}

	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", PHONENUMBER_ID)

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error al convertir a JSON:")
		return false, err
	}

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

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		//fmt.Println("✅ Mensaje enviado correctamente.")
		return true, nil
	} else {
		//fmt.Println("⚠️ Error al enviar el mensaje.")
		return false, nil
	}
}

// 2025-08-15
func dateFormat(date string) string {
	year := date[0:4]
	month := date[5:7]
	day := date[8:]
	return fmt.Sprintf("%s-%s-%s", day, month, year)
}

func totalFormat(total float64) string {
	return fmt.Sprintf("$%.2f", total)
}
