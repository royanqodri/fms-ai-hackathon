package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

	// URL API Telegram
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	}
	jsonPayload, _ := json.Marshal(payload)

	log.Println("📤 Mengirim pesan ke Telegram untuk chatID:")

	// HTTP POST
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
