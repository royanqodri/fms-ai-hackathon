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
}

type ChatServiceImpl struct {
	aiRepo operation.AiRepository
}

func NewChatService(aiRepo operation.AiRepository) ChatService {
	return &ChatServiceImpl{aiRepo: aiRepo}
}

func (service *ChatServiceImpl) AskAI(ctx *gin.Context, telegramId string, userPrompt string) (string, error) {

	// 1. Ekstrak plat nomor dari pertanyaan user
	plateNumber := ""
	words := strings.Fields(userPrompt)
	for i := 0; i < len(words)-1; i++ {
		if strings.ToUpper(words[i]) == "N" || strings.ToUpper(words[i]) == "B" || strings.ToUpper(words[i]) == "L" {
			plateNumber = words[i] + " " + words[i+1]
			break
		}
	}

	// 2. Ambil konteks data truk
	var truckContext map[string]interface{}
	if plateNumber != "" {
		truckContext, _ = service.aiRepo.GetTruckContext(ctx, nil, plateNumber)
	}

	// 3. Ambil peringatan terbaru
	recentAlerts, _ := service.aiRepo.GetRecentAlerts(ctx, nil, 5)

	// 4. Ubah konteks menjadi JSON
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

	// 7. Panggil Meta AI
	aiResponse, err := util.AskMetaAI(systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}

	return aiResponse, nil
}
