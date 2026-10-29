package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type MetaChatRequest struct {
	Model    string        `json:"model"`
	Messages []MetaMessage `json:"messages"`
}

type MetaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MetaChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// AskMetaAI for GET Meta API
func AskMetaAI(systemPrompt string, userPrompt string) (string, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		err = godotenv.Load("./cmd/config.env")
		if err != nil {
			log.Println("Warning: config.env not found, trying OS env")
		}
	}

	apiKey := os.Getenv("META_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("META_API_KEY is not set. file config.env on root folder META_API_KEY=...")
	}

	reqBody := MetaChatRequest{
		Model: "muse-spark-1.3",
		Messages: []MetaMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.meta.ai/v1/chat/completions", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Meta API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp MetaChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI")
}

func GenerateTelemetryPrompt(truckID string, temperature, weightLoad, vibration float64) (string, string) {
	systemPrompt := "You are an AI Fleet Monitoring System for Lumajang. Analyze the telemetry data and provide a brief, actionable warning if there is a risk."
	userPrompt := fmt.Sprintf("Truck ID: %s. Engine Temperature: %.2f°C. Weight Load: %.2f Tons. Vibration: %.2f.",
		truckID, temperature, weightLoad, vibration)

	return systemPrompt, userPrompt
}
