package models

type PageRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Response struct {
	Data   string   `json:"data"`
	Errors []string `json:"errors"`
}

type ResponseTest struct {
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
	Count  interface{} `json:"count,omitempty"`
}