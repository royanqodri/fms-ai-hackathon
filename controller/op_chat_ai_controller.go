package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	requestOp "github.com/royanqodri/fms-ai-hackathon/model/request/operation"
	repoOp "github.com/royanqodri/fms-ai-hackathon/repository/operation"
	serviceOp "github.com/royanqodri/fms-ai-hackathon/service/operation"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type ChatController interface {
	Chat(ctx *gin.Context)
	TelegramWebhook(ctx *gin.Context)
}

type ChatControllerImpl struct {
	chatService serviceOp.ChatService
	aiRepo      repoOp.AiRepository
}

func NewChatController(chatService serviceOp.ChatService, aiRepo repoOp.AiRepository) ChatController {
	return &ChatControllerImpl{
		chatService: chatService,
		aiRepo:      aiRepo,
	}
}

var lastPlateNumber string
var processedMessages = make(map[string]bool)
var mutex = &sync.Mutex{}

func isChartRequest(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "grafik") ||
		strings.Contains(lower, "chart") ||
		strings.Contains(lower, "graph") ||
		strings.Contains(lower, "diagram")
}

func isMuseRequest(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "muse") ||
		strings.Contains(lower, "gambar") ||
		strings.Contains(lower, "ilustrasi") ||
		strings.Contains(lower, "desain") ||
		strings.Contains(lower, "image") ||
		strings.Contains(lower, "design") ||
		strings.Contains(lower, "illustration")
}

func (controller *ChatControllerImpl) Chat(ctx *gin.Context) {
	req := requestOp.ChatRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	if isChartRequest(req.Message) {
		plateNumber := util.ExtractPlateNumber(req.Message)
		if plateNumber == "" {
			if lastPlateNumber != "" {
				plateNumber = lastPlateNumber
			} else {
				util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, "Please specify a truck plate number (e.g., 'N 5678 AB') to generate a chart.", 0, 0)
				return
			}
		} else {
			lastPlateNumber = plateNumber
		}

		if err := controller.chatService.GenerateTemperatureChart(ctx, req.TelegramId, plateNumber); err != nil {
			util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("Chart Error: %v", err), 0, 0)
			return
		}
		util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, map[string]string{"reply": fmt.Sprintf("Chart for truck %s sent to your Telegram!", plateNumber)}, "success", 0, 0)
		return
	}

	if isMuseRequest(req.Message) {
		plateNumber := util.ExtractPlateNumber(req.Message)
		if plateNumber != "" {
			lastPlateNumber = plateNumber
		} else {
			plateNumber = lastPlateNumber
		}

		// Get data from db
		truckContext, _ := controller.aiRepo.GetTruckContext(ctx, nil, plateNumber)

		prompt := "A futuristic 3D illustration of a heavy truck with engine overheating, glowing red warning lights, parked on a cyberpunk fleet management dashboard, high detail, 4k."
		if truckContext != nil {
			temp := truckContext["temperature"].(float64)
			fuel := truckContext["fuel_level"].(float64)
			oil := truckContext["oil_pressure"].(float64)
			rpm := truckContext["engine_rpm"].(float64)

			prompt = fmt.Sprintf(`A futuristic 3D bar chart design for a fleet dashboard. 
			Metrics: Temp %.1f°C (Red), Fuel %.1f%% (Blue), Oil %.1f PSI (Yellow), RPM %.0f (Cyan). 
			Style: Cyberpunk, neon, holographic, 4k.`, temp, fuel, oil, rpm)
		}

		imageData, err := util.GenerateImageWithMeta(prompt)
		if err != nil {
			util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("Muse Image Error: %v", err), 0, 0)
			return
		}

		if err := util.SendTelegramPhoto(req.TelegramId, imageData, "AI-generated illustration via Muse Image"); err != nil {
			util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("Telegram Error: %v", err), 0, 0)
			return
		}

		util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, map[string]string{"reply": "Muse Image illustration sent to your Telegram!"}, "success", 0, 0)
		return
	}

	aiReply, err := controller.chatService.AskAI(ctx, req.TelegramId, req.Message)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("AI Error: %v", err), 0, 0)
		return
	}

	err = util.SendTelegramMessage(req.TelegramId, aiReply)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("Telegram Error: %v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, map[string]string{"reply": aiReply}, "success", 0, 0)
}

func (controller *ChatControllerImpl) TelegramWebhook(ctx *gin.Context) {
	log.Println("🔔 Menerima pesan dari Telegram...")

	var update struct {
		Message struct {
			MessageID int64 `json:"message_id"`
			Chat      struct {
				ID int64 `json:"id"`
			} `json:"chat"`
			Text string `json:"text"`
		} `json:"message"`
	}

	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	mutex.Lock()
	if processedMessages[fmt.Sprintf("%d", update.Message.MessageID)] {
		mutex.Unlock()
		log.Println("⚠️ Pesan duplikat terdeteksi, diabaikan.")
		ctx.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}
	processedMessages[fmt.Sprintf("%d", update.Message.MessageID)] = true
	mutex.Unlock()

	chatID := fmt.Sprintf("%d", update.Message.Chat.ID)
	userText := update.Message.Text

	if strings.Contains(strings.ToLower(userText), "/start") {
		welcomeMsg := `🚛 *Welcome to FMS AI BOT* 🚛
_Powered by Meta AI_

Hello! I am your intelligent Fleet Management Assistant. I can help you monitor your trucks, analyze health metrics, and generate visual reports in real-time.

*Here is what I can do:*
1. 🩺 *Check Truck Health* — Ask about a specific truck's condition (e.g., "Check truck N 5678 AB").
2. 📊 *Generate Health Charts* — Visualize telemetry data (Temp, Fuel, Oil, RPM) using QuickChart.
3. 🎨 *Generate AI Illustrations* — Create futuristic dashboard designs via Meta Muse Image.
4. 🗺️ *Navigate to Nearest SPBU* — Get directions to the closest gas station if a truck is overheating.

*Try these commands:*
- "Generate a chart for truck N 5678 AB"
- "Create an illustration for truck N 5678 AB using Muse AI"
- "Where is the nearest gas station for truck N 5678 AB?"

How can I assist you today? 👇`

		util.SendTelegramMessage(chatID, welcomeMsg)
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}

	if isChartRequest(userText) {
		plateNumber := util.ExtractPlateNumber(userText)
		if plateNumber == "" {
			if lastPlateNumber != "" {
				plateNumber = lastPlateNumber
			} else {
				util.SendTelegramMessage(chatID, "Please specify a truck plate number (e.g., 'N 5678 AB') to generate a chart.")
				ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
				return
			}
		} else {
			lastPlateNumber = plateNumber
		}

		if err := controller.chatService.GenerateTemperatureChart(ctx, chatID, plateNumber); err != nil {
			util.SendTelegramMessage(chatID, fmt.Sprintf("Failed to generate chart for %s: %v", plateNumber, err))
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}

	if isMuseRequest(userText) {
		plateNumber := util.ExtractPlateNumber(userText)
		if plateNumber != "" {
			lastPlateNumber = plateNumber
		} else {
			plateNumber = lastPlateNumber
		}

		truckContext, _ := controller.aiRepo.GetTruckContext(ctx, nil, plateNumber)

		prompt := "A futuristic 3D illustration of a heavy truck with engine overheating, glowing red warning lights, parked on a cyberpunk fleet management dashboard, high detail, 4k."
		if truckContext != nil {
			temp := truckContext["temperature"].(float64)
			fuel := truckContext["fuel_level"].(float64)
			oil := truckContext["oil_pressure"].(float64)
			rpm := truckContext["engine_rpm"].(float64)

			prompt = fmt.Sprintf(`A futuristic 3D bar chart design for a fleet dashboard. 
			Metrics: Temp %.1f°C (Red), Fuel %.1f%% (Blue), Oil %.1f PSI (Yellow), RPM %.0f (Cyan). 
			Style: Cyberpunk, neon, holographic, 4k.`, temp, fuel, oil, rpm)
		}

		imageData, err := util.GenerateImageWithMeta(prompt)
		if err != nil {
			util.SendTelegramMessage(chatID, fmt.Sprintf("Muse Image Error: %v", err))
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		if err := util.SendTelegramPhoto(chatID, imageData, "AI-generated illustration via Meta Muse Image"); err != nil {
			util.SendTelegramMessage(chatID, fmt.Sprintf("Telegram Error: %v", err))
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}

	aiReply, err := controller.chatService.AskAI(ctx, chatID, userText)
	if err != nil {
		util.SendTelegramMessage(chatID, "Sorry, the AI is busy. Try again later.")
		return
	}

	util.SendTelegramMessage(chatID, aiReply)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
