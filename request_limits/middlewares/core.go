package request_limits_middlewares

import (
	"errors"
	"net/http"
)

const ResponseStatusOK = "OK"

type MiddlewareCheckRequest struct {
	Method string                     `json:"method"`
	Data   MiddlewareCheckRequestData `json:"data"`
}

type MiddlewareCheckRequestData struct {
	RequestMicroservice string `json:"request_microservice"`
	RequestMethod       string `json:"request_method"`
	ServiceStationID    string `json:"service_station_id,omitempty"`
	StoreUserID         string `json:"user_id,omitempty"`
}

type CheckResponse struct {
	Result struct {
		IsRequestAvailable bool `json:"is_request_available"`
	} `json:"result"`
	Status string `json:"status"`
}

func unauthorizedResponse(info string) (interface{}, int, error) {
	return nil, http.StatusPaymentRequired, errors.New("unauthorized:" + info)
}
