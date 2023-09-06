package responses

type ErrorResponse struct {
	ErrorType string `json:"ErrorType"`
	ErrorCode string `json:"ErrorCode"`
	Error     string `json:"Error"`
	Status    string `json:"Status"`
}

func NewErrorResponse(errorType string, error error, errorDesc string, statusCode int) (ErrorResponse, int, error) {
	return ErrorResponse{
		ErrorType: errorType,
		ErrorCode: error.Error(),
		Error:     errorDesc,
		Status:    "NOK",
	}, statusCode, error
}
