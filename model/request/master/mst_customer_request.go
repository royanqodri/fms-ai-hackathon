package master

import "time"

type MstCustomerGetRequest struct {
	Search     string    `form:"search"`
	Type       string    `form:"type"`
	StatusData string    `form:"status_data"`
	Source     string    `form:"source" binding:"required"`
	LastTime   time.Time `form:"last_time" time_format:"2006-01-02 15:04:05"`
}

type MstCustomerPostRequest struct {
	Source string                     `json:"source" binding:"required"`
	Data   []MstCustomerDetailRequest `json:"data" binding:"required,dive"`
}

type MstCustomerDetailRequest struct {
	Id         string `json:"id"`
	CustomerNo string `json:"customer_no" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	PhoneNo    string `json:"phone_no"`
	InsertBy   string `json:"insert_by"`
	UpdateBy   string `json:"update_by"`
}
