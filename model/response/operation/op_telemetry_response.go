package operation

type TrxTelemetryGetResponse struct {
	Id          string  `json:"id"`
	TruckId     string  `json:"truck_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WeightLoad  float64 `json:"weight_load"`
	Vibration   float64 `json:"vibration"`
	FuelLevel   float64 `json:"fuel_level"`
	OilPressure float64 `json:"oil_pressure"`
	EngineRpm   float64 `json:"engine_rpm"`
	RecordedAt  string  `json:"recorded_at"`
}

type TrxTelemetryPostResponse struct {
	Id         string `json:"id"`
	TruckId    string `json:"truck_id"`
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}
