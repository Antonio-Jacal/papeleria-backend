package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/resend/resend-go/v3"
)

func SendHTMLEmail(to []string, subject string, htmlBody string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("EMAIL_SMTP")
	smtpPort := os.Getenv("EMAIL_PORT")

	if smtpPort == "" {
		smtpPort = "587"
	}

	// Encabezados y cuerpo
	message := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		htmlBody)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("error al enviar correo: %w", err)
	}

	return nil
}

func SendHTMLEmailResend(to []string, subject string, htmlBody string) error {
	API_KEY_RESEND := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(API_KEY_RESEND)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"antoniojacalm@gmail.com"},
		Subject: subject,
		Html:    htmlBody,
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Fatalf("failed to send email: %v", err)
	}

	log.Printf("email sent: %s", sent.Id)
	return nil
}

func FormatDate(t *time.Time) string {
	if t == nil {
		return "N/A"
	}
	return t.Format("02 Jan 2006")
}

func BuildProductosHTML(productos map[string]models.ProductoDetalle) string {
	var sb strings.Builder
	sb.WriteString("<ul style='padding-left: 20px;'>")
	for nombre, detalle := range productos {
		if detalle.Cantidad > 0 {
			sb.WriteString(fmt.Sprintf("<li>%s: %d</li>", nombre, detalle.Cantidad))
		}
	}
	sb.WriteString("</ul>")
	return sb.String()
}

func BuildUtilesQuitadosHTML(utiles map[string]int) string {
	if len(utiles) == 0 {
		return "<p>No se eliminaron útiles.</p>"
	}
	var sb strings.Builder
	sb.WriteString("<ul style='padding-left: 20px;'>")
	for nombre, cantidad := range utiles {
		if cantidad > 0 {
			sb.WriteString(fmt.Sprintf("<li>%s: %d</li>", nombre, cantidad))
		}
	}
	sb.WriteString("</ul>")
	return sb.String()
}
