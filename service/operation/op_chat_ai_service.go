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

	// Prompt
	systemPrompt := fmt.Sprintf(`
		You are an intelligent AI Assistant for the Fleet Management System (FMS).
		Your task is to answer the fleet manager's questions accurately, concisely, and professionally based ONLY on the provided context data.
		
		CONTEXT DATA:
		%s
		
		INSTRUCTIONS:
		1. If the user asks about a specific truck, use the "requested_truck_data" if available.
		2. If the "requested_truck_data" is null or empty, politely inform the user that the specific truck data is not found in the database.
		3. If the user asks about general fleet issues, refer to "recent_fleet_alerts".
		4. Do not hallucinate or invent data. If the information is not in the context, state that you don't know.
		5. Keep your answer brief and actionable.
	`, string(contextData))

	// GET Meta AI
	aiResponse, err := util.AskMetaAI(systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}

	return aiResponse, nil
}
