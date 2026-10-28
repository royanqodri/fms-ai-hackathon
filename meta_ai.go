package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Struct untuk request ke Meta API
type MetaChatRequest struct {
	Model    string        `json:"model"`
	Messages []MetaMessage `json:"messages"`
}

type MetaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Struct untuk response dari Meta API
type MetaChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// AskMetaAI for GET Meta MAPI
func AskMetaAI(systemPrompt string, userPrompt string) (string, error) {
	apiKey := os.Getenv("META_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("META_API_KEY is not set in .env")
	}

	reqBody := MetaChatRequest{
		Model: "muse-spark-1.3", // Gunakan model terbaru Meta
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
