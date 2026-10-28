package operation

type ChatResponse struct {
	TelegramId string `json:"telegram_id"`
	Reply      string `json:"reply"`
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}
