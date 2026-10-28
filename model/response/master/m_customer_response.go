package master

type MstCustomerGetResponse struct {
	Id         string `json:"id"`
	CustomerNo string `json:"customer_no"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	PhoneNo    string `json:"phone_no"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type MstCustomerPostResponse struct {
	Id         string `json:"id"`
	CustomerNo string `json:"customer_no"`
	Name       string `json:"name"`
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}
