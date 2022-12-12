package handler

type Response struct {
	Data   string   `json:"data"`
	Errors []string `json:"errors"`
}