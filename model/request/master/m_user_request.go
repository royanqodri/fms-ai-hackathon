package master

import "time"

type MstUserGetRequest struct {
	Search     string    `form:"search"`
	Role       string    `form:"role"`
	StatusData string    `form:"status_data"`
	Source     string    `form:"source" binding:"required"`
	LastTime   time.Time `form:"last_time" time_format:"2006-01-02 15:04:05"`
}

type MstUserPostRequest struct {
	Source string                 `json:"source" binding:"required"`
	Data   []MstUserDetailRequest `json:"data" binding:"required,dive"`
}

type MstUserDetailRequest struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id" binding:"required,uuid"`
	TelegramId string `json:"telegram_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Role       string `json:"role" binding:"required"`
	InsertBy   string `json:"insert_by"`
	UpdateBy   string `json:"update_by"`
}
