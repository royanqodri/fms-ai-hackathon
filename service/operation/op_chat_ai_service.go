package operation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/fms-ai-hackathon/repository/operation"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type ChatService interface {
	AskAI(ctx *gin.Context, telegramId string, userPrompt string) (string, error)
	GenerateTemperatureChart(ctx *gin.Context, chatID string, plateNumber string) error
}

type ChatServiceImpl struct {
	aiRepo operation.AiRepository
}

func NewChatService(aiRepo operation.AiRepository) ChatService {
	return &ChatServiceImpl{aiRepo: aiRepo}
}

func (service *ChatServiceImpl) AskAI(ctx *gin.Context, telegramId string, userPrompt string) (string, error) {

	plateNumber := ""
	words := strings.Fields(userPrompt)
	for i := 0; i < len(words)-1; i++ {
		if strings.ToUpper(words[i]) == "N" || strings.ToUpper(words[i]) == "B" || strings.ToUpper(words[i]) == "L" {
			plateNumber = words[i] + " " + words[i+1]
			break
		}
	}

	var truckContext map[string]interface{}
	if plateNumber != "" {
		truckContext, _ = service.aiRepo.GetTruckContext(ctx, nil, plateNumber)
	}

	recentAlerts, _ := service.aiRepo.GetRecentAlerts(ctx, nil, 5)

	contextData, _ := json.Marshal(map[string]interface{}{
		"requested_truck_data": truckContext,
		"recent_fleet_alerts":  recentAlerts,
	})

	spbuList := `
	1. SPBU Lumajang Kota: Lat -8.133, Lng 113.223
	2. SPBU Candipuro: Lat -8.190, Lng 113.075
	3. SPBU Senduro: Lat -8.100, Lng 113.050
	4. SPBU Pasirian: Lat -8.200, Lng 113.100
	5. SPBU Tempeh: Lat -8.180, Lng 113.150
	6. SPBU Kunir: Lat -8.220, Lng 113.180
	`

	systemPrompt := fmt.Sprintf(`
		You are an intelligent AI Assistant for the Fleet Management System (FMS) in Lumajang.
		Your task is to answer the fleet manager's questions accurately, concisely, and professionally based ONLY on the provided context data.
		
		CONTEXT DATA (Truck & Alerts):
		%s
		
		LIST OF GAS STATIONS (SPBU) IN LUMAJANG:
		%s
		
		INSTRUCTIONS:
		1. If the user asks about a specific truck, use the "requested_truck_data" if available.
		2. If the "requested_truck_data" is null or empty, politely inform the user that the specific truck data is not found in the database.
		3. If the user asks for directions to the nearest gas station (SPBU), calculate which SPBU from the list is closest to the truck's current Latitude and Longitude, and provide the name and directions.
		4. If the user asks about general fleet issues, refer to "recent_fleet_alerts".
		5. Do not hallucinate or invent data. If the information is not in the context, state that you don't know.
		6. Keep your answer brief and actionable.
	`, string(contextData), spbuList)

	// Call AI
	aiResponse, err := util.AskMetaAI(systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}

	return aiResponse, nil
}

func (service *ChatServiceImpl) GenerateTemperatureChart(ctx *gin.Context, chatID string, plateNumber string) error {
	history, err := service.aiRepo.GetTemperatureHistory(ctx, nil, plateNumber, 1)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}

	if len(history) == 0 {
		return fmt.Errorf("no telemetry data found for truck %s", plateNumber)
	}

	var labels []string
	var tempData, fuelData, oilData, rpmData []float64

	for _, h := range history {
		timeStr := h["recorded_at"].(string)
		if len(timeStr) >= 16 {
			labels = append(labels, timeStr[11:16])
		} else {
			labels = append(labels, timeStr)
		}
		tempData = append(tempData, h["temperature"].(float64))
		fuelData = append(fuelData, h["fuel_level"].(float64))
		oilData = append(oilData, h["oil_pressure"].(float64))
		rpmData = append(rpmData, h["engine_rpm"].(float64)/100)
	}

	title := fmt.Sprintf("Latest Health Metrics - Truck %s", plateNumber)
	return util.SendChartToTelegram(chatID, title, labels, tempData, fuelData, oilData, rpmData)
}
