package operation

import "time"

type LogAlertGetRequest struct {
	TruckId    string    `form:"truck_id"`
	AlertType  string    `form:"alert_type"`
	IsResolved *bool     `form:"is_resolved"`
	Source     string    `form:"source" binding:"required"`
	LastTime   time.Time `form:"last_time" time_format:"2006-01-02 15:04:05"`
}

type LogAlertPostRequest struct {
	Source string                  `json:"source" binding:"required"`
	Data   []LogAlertDetailRequest `json:"data" binding:"required,dive"`
}

type LogAlertDetailRequest struct {
	Id         string `json:"id"`
	TruckId    string `json:"truck_id" binding:"required,uuid"`
	AlertType  string `json:"alert_type" binding:"required"`
	AiAnalysis string `json:"ai_analysis"`
	IsResolved bool   `json:"is_resolved"`
	UpdateBy   string `json:"update_by"`
}
