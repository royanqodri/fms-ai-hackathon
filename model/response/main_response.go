package response

type StatusResponse struct {
	Code    int    `json:"status_code"`
	Status  string `json:"status_text"`
	Message string `json:"message"`
}

type MainResponse struct {
	StatusResponse StatusResponse `json:"status_response"`
	TotalPage      int64          `json:"total_page"`
	TotalData      int64          `json:"total_data"`
	DataHeader     interface{}    `json:"data_header"`
	Data           interface{}    `json:"data"`
}
