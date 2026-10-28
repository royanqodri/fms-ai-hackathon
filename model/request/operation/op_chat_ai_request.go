package operation

type ChatRequest struct {
	TelegramId string `json:"telegram_id" binding:"required"`
	Message    string `json:"message" binding:"required"`
}
