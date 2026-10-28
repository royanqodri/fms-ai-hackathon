package master

type MstUserGetResponse struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id"`
	TelegramId string `json:"telegram_id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type MstUserPostResponse struct {
	Id         string `json:"id"`
	TelegramId string `json:"telegram_id"`
	Name       string `json:"name"`
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}
