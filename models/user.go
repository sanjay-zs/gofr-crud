package models

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Response struct {
	User       interface{} `json:"product"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status-code"`
}
