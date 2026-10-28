package operation

type LogAlertGetResponse struct {
	Id         string `json:"id"`
	TruckId    string `json:"truck_id"`
	AlertType  string `json:"alert_type"`
	AiAnalysis string `json:"ai_analysis"`
	IsResolved bool   `json:"is_resolved"`
	CreatedAt  string `json:"created_at"`
}

type LogAlertPostResponse struct {
	Id         string `json:"id"`
	TruckId    string `json:"truck_id"`
	AlertType  string `json:"alert_type"`
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}
