package master

import "time"

type MstTruckGetRequest struct {
	Search     string    `form:"search"`
	CargoType  string    `form:"cargo_type"`
	Status     string    `form:"status"`
	CustomerId string    `form:"customer_id"`
	StatusData string    `form:"status_data"`
	Source     string    `form:"source" binding:"required"`
	LastTime   time.Time `form:"last_time" time_format:"2006-01-02 15:04:05"`
}

type MstTruckPostRequest struct {
	Source string                  `json:"source" binding:"required"`
	Data   []MstTruckDetailRequest `json:"data" binding:"required,dive"`
}

type MstTruckDetailRequest struct {
	Id          string `json:"id"`
	CustomerId  string `json:"customer_id" binding:"required,uuid"`
	PlateNumber string `json:"plate_number" binding:"required"`
	CargoType   string `json:"cargo_type" binding:"required"`
	Status      string `json:"status" binding:"required"`
	InsertBy    string `json:"insert_by"`
	UpdateBy    string `json:"update_by"`
}
