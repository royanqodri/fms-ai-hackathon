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

	// list koordinate real fuel in lumajang base
	spbuList := `
	No | Kode SPBU | Nama & Lokasi | Latitude | Longitude | Jalur
	1  | 54.673.01 | SPBU Pasirian, Jl. Raya Pasirian - Lumajang | -8.2214 | 113.1200 | MAIN_ROAD
	2  | 54.673.02 | SPBU Sukodono, Jl. Gatot Subroto No.3, Karangsari | -8.1125 | 113.2294 | MAIN_ROAD
    3  | 54.673.13 | SPBU Labruk, Jl. Slamet Wardoyo 103, Ds. Labruk Lor | -8.1394 | 113.2086 | MAIN_ROAD
	4  | 53.673.17 | SPBU Condro Pasirian, Jl. Nasional III, Kel. Condro | -8.2050 | 113.1450 | MAIN_ROAD
	5  | 54.673.07 | SPBU Tempeh, Jl. Raya Tempeh KM 8 | -8.180560 | 113.181637 | MAIN_ROAD
	`

	systemPrompt := fmt.Sprintf(`
		You are an intelligent AI Assistant for the Fleet Management System (FMS) in Lumajang.
		Your task is to answer the fleet manager's questions accurately, concisely, and professionally based ONLY on the provided context data.
		
		CONTEXT DATA (Truck & Alerts):
		%s
		
		LIST OF GAS STATIONS (SPBU) IN LUMAJANG (ONLY MAIN ROAD / JALUR MERAH):
		%s
		
		INSTRUCTIONS:
		1. If the user asks about a specific truck, use the "requested_truck_data" if available.
		2. If the "requested_truck_data" is null or empty, politely inform the user that the specific truck data is not found in the database.
		3. If the user asks for directions to the nearest gas station (SPBU), you MUST ONLY choose from the list provided above (which are all on the main road). Calculate the nearest one based on the truck's Latitude and Longitude. Provide the SPBU code, name, and direction.
		4. Do NOT recommend any SPBU outside of this list (to avoid routes that go off the main road).
		5. If the user asks about general fleet issues, refer to "recent_fleet_alerts".
		6. Do not hallucinate or invent data. If the information is not in the context, state that you don't know.
		7. Keep your answer brief and actionable.
		8. Respond in the same language the user is using.
	`, string(contextData), spbuList)

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
