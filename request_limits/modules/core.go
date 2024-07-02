package request_limits_modules

import "go.mongodb.org/mongo-driver/bson"

const ResponseStatusOK = "OK"

type MicroserviceDecrementRequest struct {
	Method   string                           `json:"method"`
	Metadata bson.M                           `json:"metadata"`
	Data     DecrementRequestLimitRequestData `json:"data"`
}

type MicroserviceIncrementRequest struct {
	Method   string                           `json:"method"`
	Metadata bson.M                           `json:"metadata"`
	Data     IncrementRequestLimitRequestData `json:"data"`
}

type DecrementRequestLimitRequestData struct {
	RequestMicroservice string `json:"request_microservice"`
	RequestMethod       string `json:"request_method"`
	ServiceStationID    string `json:"service_station_id,omitempty"`
	StoreUserID         string `json:"user_id,omitempty"`
}

type IncrementRequestLimitRequestData struct {
	RequestMicroservice string `json:"request_microservice"`
	RequestMethod       string `json:"request_method"`
	ServiceStationID    string `json:"service_station_id,omitempty"`
	StoreUserID         string `json:"user_id,omitempty"`
	Date                int64  `json:"date"`
}
type MicroserviceResponse struct {
	ErrorType string `json:"ErrorType"`
	ErrorCode string `json:"ErrorCode"`
	Error     string `json:"Error"`
	Status    string `json:"Status"`
}
