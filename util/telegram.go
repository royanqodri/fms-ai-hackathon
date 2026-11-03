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

	log.Println("📤 Sending message to Telegram for chatID:", chatID)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("❌ Failed to POST message to Telegram:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("❌ Failed to send Telegram message, status:", resp.StatusCode)
		return fmt.Errorf("failed to send telegram message, status: %d", resp.StatusCode)
	}

	log.Println("✅ Message successfully sent to Telegram!")
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

	log.Println("📤 Uploading photo to Telegram for chatID:", chatID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("❌ Failed to POST photo to Telegram:", err)
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Failed to send photo to Telegram, status: %d, body: %s\n", resp.StatusCode, string(respBody))
		return fmt.Errorf("failed to send photo, status: %d", resp.StatusCode)
	}

	log.Println("✅ Photo successfully sent to Telegram!")
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

func DownloadTelegramFile(fileID string) (string, error) {

	err := godotenv.Load("config.env")
	if err != nil {
		_ = godotenv.Load("./cmd/config.env")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return "", fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}

	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", botToken, fileID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Ok     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if !result.Ok {
		return "", fmt.Errorf("failed to get file path from Telegram")
	}

	fileUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", botToken, result.Result.FilePath)
	respFile, err := http.Get(fileUrl)
	if err != nil {
		return "", err
	}
	defer respFile.Body.Close()

	if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
		_ = os.Mkdir("./tmp", 0755)
	}
	localPath := fmt.Sprintf("./tmp/%s.ogg", fileID)
	out, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, respFile.Body)
	if err != nil {
		return "", err
	}

	log.Println("✅ Audio file downloaded from Telegram!")
	return localPath, nil
}
