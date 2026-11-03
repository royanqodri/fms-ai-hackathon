package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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

func GenerateTelemetryPrompt(truckID string, latitude, longitude, temperature, weightLoad, vibration, fuelLevel, oilPressure, engineRpm float64) (string, string) {
	systemPrompt := `You are an AI Fleet Monitoring System for Lumajang. 
Analyze the telemetry data and provide a brief, actionable warning if there is a risk. 
If the engine is overheating, suggest the driver to pull over safely and find the nearest gas station (SPBU) or rest area based on their current coordinates.`

	userPrompt := fmt.Sprintf(
		"Truck ID: %s. Current Location: Latitude %.6f, Longitude %.6f. "+
			"Engine Temperature: %.2f°C. Weight Load: %.2f Tons. Vibration: %.2f. "+
			"Fuel Level: %.2f%%. Oil Pressure: %.2f PSI. Engine RPM: %.2f.",
		truckID, latitude, longitude, temperature, weightLoad, vibration, fuelLevel, oilPressure, engineRpm,
	)

	return systemPrompt, userPrompt
}

// GenerateImageWithMeta Call Meta Muse Image
func GenerateImageWithMeta(prompt string) ([]byte, error) {
	// Load API Key
	err := godotenv.Load("config.env")
	if err != nil {
		_ = godotenv.Load("./cmd/config.env")
	}

	apiKey := os.Getenv("META_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("META_API_KEY is not set")
	}

	url := "https://api.meta.ai/v1/images/generations"

	payload := map[string]interface{}{
		"model":           "muse-image-1.0",
		"prompt":          prompt,
		"n":               1,
		"size":            "1024x1024",
		"response_format": "b64_json",
	}
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Meta Image API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse b64_json
	var result struct {
		Data []struct {
			B64Json string `json:"b64_json"`
		} `json:"data"`
	}
	json.Unmarshal(body, &result)

	if len(result.Data) > 0 && result.Data[0].B64Json != "" {
		imgData, err := base64.StdEncoding.DecodeString(result.Data[0].B64Json)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 image: %w", err)
		}
		log.Println("✅ Ilustration AI success!")
		return imgData, nil
	}

	return nil, fmt.Errorf("no image data returned. Raw response: %s", string(body))
}

func TranscribeAudioWithMeta(audioFilePath string) (string, error) {
	// Load API Key
	_ = godotenv.Load("config.env")
	_ = godotenv.Load("./cmd/config.env")

	apiKey := os.Getenv("META_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("META_API_KEY is not set")
	}

	wavPath := filepath.Join(filepath.Dir(audioFilePath), fmt.Sprintf("converted_%d.wav", time.Now().UnixNano()))

	ffmpegPath := "./ffmpeg.exe"

	cmd := exec.Command(ffmpegPath, "-y", "-i", audioFilePath, "-ac", "1", "-ar", "24000", "-c:a", "pcm_s16le", "-map_metadata", "-1", wavPath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to convert audio to WAV: %w | stderr: %s", err, stderr.String())
	}
	defer os.Remove(wavPath)

	audioFile, err := os.Open(wavPath)
	if err != nil {
		return "", fmt.Errorf("failed to open WAV file: %w", err)
	}
	defer audioFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="request"`)
	h.Set("Content-Type", "application/json")

	part, err := writer.CreatePart(h)
	if err != nil {
		return "", fmt.Errorf("failed to create form part 'request': %w", err)
	}

	requestJson := map[string]interface{}{
		"model":         "muse-voice-transcribe-1.0",
		"audioEncoding": "WAV",
		"mode":          "PUSH_TO_TALK",
	}
	requestBytes, _ := json.Marshal(requestJson)
	_, _ = part.Write(requestBytes)

	audioPart, err := writer.CreateFormFile("audio", "audio.wav")
	if err != nil {
		return "", fmt.Errorf("failed to create form file 'audio': %w", err)
	}
	_, _ = io.Copy(audioPart, audioFile)
	_ = writer.Close()

	url := "https://api.meta.ai/v1/asr/transcribe?sessionId=fms-hackathon-session"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Meta API: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Meta Voice API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Transcript string `json:"transcript"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	log.Println("✅ Audio transcribed successfully!")
	return result.Transcript, nil
}
