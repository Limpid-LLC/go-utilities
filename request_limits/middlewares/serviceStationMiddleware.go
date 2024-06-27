package request_limits_middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Limpid-LLC/saiService"
	"io/ioutil"
	"net/http"
)

func CreateServiceStationRequestLimitMiddleware(requestLimitServiceURL string, microserviceName string, method string) func(next saiService.HandlerFunc, data interface{}, metadata interface{}) (interface{}, int, error) {
	return func(next saiService.HandlerFunc, data interface{}, metadata interface{}) (interface{}, int, error) {
		if requestLimitServiceURL == "" {
			fmt.Println("serviceStationRequestLimitMiddleware: requestLimit service url is empty")
			return unauthorizedResponse("requestLimitServiceURL")
		}

		var dataMap map[string]interface{}

		dataBytes, _ := json.Marshal(data)

		_ = json.Unmarshal(dataBytes, &dataMap)

		var serviceStationID string
		if dataMap["service_station_id"] != nil {
			serviceStationID = dataMap["service_station_id"].(string)
		} else if dataMap["sto_id"] != nil {
			serviceStationID = dataMap["sto_id"].(string)
		} else if dataMap["stoId"] != nil {
			serviceStationID = dataMap["stoId"].(string)
		} else {
			serviceStationID = ""
		}

		if len(serviceStationID) == 0 {
			return unauthorizedResponse("empty service_station_id")
		}

		checkReq := MiddlewareCheckRequest{
			Method: "check",
			Data: MiddlewareCheckRequestData{
				Microservice:     microserviceName,
				Method:           method,
				ServiceStationID: serviceStationID,
			},
		}

		jsonData, err := json.Marshal(checkReq)
		if err != nil {
			fmt.Println("serviceStationRequestLimitMiddleware: error marshaling data")
			fmt.Println("serviceStationRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("marshaling -> " + err.Error())
		}

		fmt.Println("serviceStationRequestLimitMiddleware checkRequest: " + string(jsonData))

		req, err := http.NewRequest("POST", requestLimitServiceURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("serviceStationRequestLimitMiddleware: error creating request")
			fmt.Println("serviceStationRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("creating request -> " + err.Error())
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("serviceStationRequestLimitMiddleware: error sending request to auth")
			fmt.Println("serviceStationRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("sending request -> " + err.Error())
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		fmt.Println("serviceStationRequestLimitMiddleware Response Status: " + resp.Status)
		fmt.Println("serviceStationRequestLimitMiddleware Response Body: " + string(body))

		if err != nil {
			fmt.Println("serviceStationRequestLimitMiddleware: error reading body from auth")
			fmt.Println("serviceStationRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("reading body -> " + err.Error())
		}

		var response CheckResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("serviceStationRequestLimitMiddleware: error unmarshalling body from auth")
			fmt.Println("serviceStationRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("Unmarshal -> " + err.Error())
		}

		fmt.Println("serviceStationRequestLimitMiddleware checkResponse Body Status: " + response.Status)

		if response.Status != ResponseStatusOK {
			fmt.Println("serviceStationRequestLimitMiddleware: response-body -> result is not `Ok`")
			fmt.Println("serviceStationRequestLimitMiddleware: " + string(body))
			return unauthorizedResponse("Result -> " + string(body))
		} else {
			fmt.Println("serviceStationRequestLimitMiddleware checkResponse result check: " + fmt.Sprintf("%v", response.Result.IsRequestAvailable))
			if response.Result.IsRequestAvailable {
				return next(data, metadata)
			} else {
				return unauthorizedResponse("request limit exceeded")
			}
		}
	}
}
