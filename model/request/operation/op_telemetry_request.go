package operation

import "time"

type TrxTelemetryGetRequest struct {
	TruckId   string    `form:"truck_id"`
	StartDate string    `form:"start_date" time_format:"2006-01-02"`
	EndDate   string    `form:"end_date" time_format:"2006-01-02"`
	Source    string    `form:"source" binding:"required"`
	LastTime  time.Time `form:"last_time" time_format:"2006-01-02 15:04:05"`
}

type TrxTelemetryPostRequest struct {
	Source string                      `json:"source" binding:"required"`
	Data   []TrxTelemetryDetailRequest `json:"data" binding:"required,dive"`
}

type TrxTelemetryDetailRequest struct {
	Id          string  `json:"id"`
	TruckId     string  `json:"truck_id" binding:"required,uuid"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Temperature float64 `json:"temperature" binding:"required"`
	WeightLoad  float64 `json:"weight_load"`
	Vibration   float64 `json:"vibration"`
	FuelLevel   float64 `json:"fuel_level"`
	OilPressure float64 `json:"oil_pressure"`
	EngineRpm   float64 `json:"engine_rpm"`
	InsertBy    string  `json:"insert_by"`
}
