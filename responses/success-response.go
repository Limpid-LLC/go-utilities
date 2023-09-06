package responses

import "net/http"

type SuccessResponse struct {
	Result interface{} `json:"result,omitempty"`
	Status string      `json:"status"`
}

func NewOkResponse(result interface{}) (SuccessResponse, int, error) {
	return SuccessResponse{
		Result: result,
		Status: "OK",
	}, http.StatusOK, nil
}
