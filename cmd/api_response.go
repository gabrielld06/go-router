package router

type ApiResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}
