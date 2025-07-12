package utils

import (
	"bytes"
	"fmt"

	"github.com/Antonio-Jacal/papeleria-backend.git/models"
	"github.com/jung-kurt/gofpdf"
)

func GenerarPDFLista(productos map[string]models.ProductoDetalle) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Lista de utiles escolares", false)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 10, "Lista de útiles escolares")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)

	for nombre, detalle := range productos {
		linea := pdf.UnicodeTranslatorFromDescriptor("")(fmt.Sprintf("- %s: %d", nombre, detalle.Cantidad))
		pdf.MultiCell(0, 8, linea, "", "", false)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
