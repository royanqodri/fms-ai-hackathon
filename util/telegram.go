package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendTelegramMessage(chatID string, message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is not set in .env")
	}

	// URL API Telegram
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	}
	jsonPayload, _ := json.Marshal(payload)

	// send HTTP POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send telegram message, status: %d", resp.StatusCode)
	}

	return nil
}
