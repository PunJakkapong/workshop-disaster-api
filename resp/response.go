package resp

// ErrorResponse ..
type ErrorResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse ..
type SuccessResponse struct {
	Code    uint        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
