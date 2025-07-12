package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func UploadPDFToSupabase(pdf []byte, fileName string) (string, error) {
	supabaseUrl := os.Getenv("SUPABASE_URL")     // ej. https://xyzcompany.supabase.co
	supabaseKey := os.Getenv("SUPABASE_API_KEY") // clave anon
	bucket := os.Getenv("SUPABASE_BUCKET")       // ej. "documentos"

	uploadUrl := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseUrl, bucket, fileName)

	req, err := http.NewRequest("POST", uploadUrl, bytes.NewReader(pdf))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", "application/pdf")
	req.Header.Set("x-upsert", "true")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Construimos el link público manualmente:
		publicUrl := fmt.Sprintf("%s/storage/v1/object/public/%s//%s", supabaseUrl, bucket, fileName)
		return publicUrl, nil
	}

	body, _ := io.ReadAll(resp.Body)
	return "", fmt.Errorf("error al subir a supabase: %s", body)
}
