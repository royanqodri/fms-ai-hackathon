package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	requestOp "github.com/royanqodri/fms-ai-hackathon/model/request/operation"
	serviceOp "github.com/royanqodri/fms-ai-hackathon/service/operation"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type ChatController interface {
	Chat(ctx *gin.Context)
	TelegramWebhook(ctx *gin.Context)
}

type ChatControllerImpl struct {
	chatService serviceOp.ChatService
}

func NewChatController(chatService serviceOp.ChatService) ChatController {
	return &ChatControllerImpl{chatService: chatService}
}

func (controller *ChatControllerImpl) Chat(ctx *gin.Context) {
	req := requestOp.ChatRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
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
			Chat struct {
				ID int64 `json:"id"`
			} `json:"chat"`
			Text string `json:"text"`
		} `json:"message"`
	}

	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	chatID := fmt.Sprintf("%d", update.Message.Chat.ID)
	aiReply, err := controller.chatService.AskAI(ctx, chatID, update.Message.Text)
	if err != nil {

		util.SendTelegramMessage(chatID, "sorry, The AI ​​is busy. Try again later..")
		return
	}

	// Send AI reply to Telegram
	util.SendTelegramMessage(chatID, aiReply)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
