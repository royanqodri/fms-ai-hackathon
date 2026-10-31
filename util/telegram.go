package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func SendTelegramMessage(chatID string, message string) error {
	// load config.env
	err := godotenv.Load("config.env")
	if err != nil {
		err = godotenv.Load("./cmd/config.env")
		if err != nil {
			log.Println("⚠️ Warning: config.env not found, trying OS env")
		}
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Println("❌ TELEGRAM_BOT_TOKEN is not set in config.env")
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is not set in config.env")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	}
	jsonPayload, _ := json.Marshal(payload)

	log.Println("📤 Mengirim pesan ke Telegram untuk chatID:", chatID)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("❌ Gagal HTTP POST ke Telegram:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("❌ Gagal kirim Telegram, status:", resp.StatusCode)
		return fmt.Errorf("failed to send telegram message, status: %d", resp.StatusCode)
	}

	log.Println("✅ Pesan berhasil dikirim ke Telegram!")
	return nil
}

func SendTelegramPhoto(chatID string, imageData []byte, caption string) error {
	// Load config.env
	err := godotenv.Load("config.env")
	if err != nil {
		err = godotenv.Load("./cmd/config.env")
		if err != nil {
			log.Println("⚠️ Warning: config.env not found, trying OS env")
		}
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Println("❌ TELEGRAM_BOT_TOKEN is not set in config.env")
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is not set in config.env")
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("chat_id", chatID)
	if caption != "" {
		_ = writer.WriteField("caption", caption)
	}

	part, err := writer.CreateFormFile("photo", "chart.png")
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = part.Write(imageData)
	if err != nil {
		return fmt.Errorf("failed to write image data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// send POST request to Telegram
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", botToken)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Println("📤 Mengunggah foto ke Telegram untuk chatID:", chatID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("❌ Gagal HTTP POST foto ke Telegram:", err)
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Gagal kirim foto Telegram, status: %d, body: %s\n", resp.StatusCode, string(respBody))
		return fmt.Errorf("failed to send photo, status: %d", resp.StatusCode)
	}

	log.Println("✅ Foto berhasil dikirim ke Telegram!")
	return nil
}

func SendChartToTelegram(chatID string, title string, labels []string, tempData, fuelData, oilData, rpmData []float64) error {
	imageData, err := GenerateChartImage(title, labels, tempData, fuelData, oilData, rpmData)
	if err != nil {
		return fmt.Errorf("failed to generate chart: %w", err)
	}

	err = SendTelegramPhoto(chatID, imageData, fmt.Sprintf("Here is the health chart: %s", title))
	if err != nil {
		return fmt.Errorf("failed to send chart to telegram: %w", err)
	}

	return nil
}
