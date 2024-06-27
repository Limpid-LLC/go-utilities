package request_limits_middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/Limpid-LLC/saiService"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateStoreRequestLimitMiddleware(requestLimitServiceURL string, microserviceName string, method string) func(next saiService.HandlerFunc, data interface{}, metadata interface{}) (interface{}, int, error) {
	return func(next saiService.HandlerFunc, data interface{}, metadata interface{}) (interface{}, int, error) {
		if requestLimitServiceURL == "" {
			log.Println("storeRequestLimitMiddleware: requestLimit service url is empty")
			return unauthorizedResponse("requestLimitServiceURL")
		}

		var dataMap map[string]interface{}

		dataBytes, _ := json.Marshal(data)

		_ = json.Unmarshal(dataBytes, &dataMap)

		checkReq := MiddlewareCheckRequest{
			Method: "check",
			Data: MiddlewareCheckRequestData{
				Microservice: microserviceName,
				Method:       method,
				StoreUserID:  dataMap["user_id"].(string),
			},
		}

		jsonData, err := json.Marshal(checkReq)
		if err != nil {
			log.Println("storeRequestLimitMiddleware: error marshaling data")
			log.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("marshaling -> " + err.Error())
		}

		req, err := http.NewRequest("POST", requestLimitServiceURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("storeRequestLimitMiddleware: error creating request")
			log.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("creating request -> " + err.Error())
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("storeRequestLimitMiddleware: error sending request to auth")
			log.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("sending request -> " + err.Error())
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Println("storeRequestLimitMiddleware: error reading body from auth")
			log.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("reading body -> " + err.Error())
		}

		var response CheckResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("storeRequestLimitMiddleware: error unmarshalling body from auth")
			log.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("Unmarshal -> " + err.Error())
		}

		if response.Status != ResponseStatusOK {
			log.Println("storeRequestLimitMiddleware: response-body -> result is not `Ok`")
			log.Println("storeRequestLimitMiddleware: " + string(body))
			return unauthorizedResponse("Result -> " + string(body))
		} else if !response.Result.IsRequestAvailable {
			return unauthorizedResponse(" request limit exceeded")
		}

		return next(data, metadata)
	}
}
