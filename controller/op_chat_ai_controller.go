package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	requestOp "github.com/royanqodri/fms-ai-hackathon/model/request/operation"
	serviceOp "github.com/royanqodri/fms-ai-hackathon/service/operation"
	"github.com/royanqodri/fms-ai-hackathon/util"
)

type ChatController interface {
	Chat(ctx *gin.Context)
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

	// Get Service AI
	aiReply, err := controller.chatService.AskAI(ctx, req.TelegramId, req.Message)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, map[string]string{"reply": aiReply}, "success", 0, 0)
}
