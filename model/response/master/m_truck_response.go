package master

type MstTruckGetResponse struct {
	Id          string `json:"id"`
	CustomerId  string `json:"customer_id"`
	PlateNumber string `json:"plate_number"`
	CargoType   string `json:"cargo_type"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type MstTruckPostResponse struct {
	Id          string `json:"id"`
	PlateNumber string `json:"plate_number"`
	CargoType   string `json:"cargo_type"`
	StatusCode  int    `json:"status_code"`
	StatusText  string `json:"status_text"`
	Message     string `json:"message"`
}
