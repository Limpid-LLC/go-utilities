package request_limits_middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Limpid-LLC/saiService"
	"io/ioutil"
	"net/http"
)

func CreateStoreRequestLimitMiddleware(requestLimitServiceURL string, microserviceName string, method string) func(next saiService.HandlerFunc, data interface{}, metadata interface{}) (interface{}, int, error) {
	return func(next saiService.HandlerFunc, data interface{}, metadata interface{}) (interface{}, int, error) {
		if requestLimitServiceURL == "" {
			fmt.Println("storeRequestLimitMiddleware: requestLimit service url is empty")
			return unauthorizedResponse("requestLimitServiceURL")
		}

		var dataMap map[string]interface{}

		dataBytes, _ := json.Marshal(data)

		_ = json.Unmarshal(dataBytes, &dataMap)

		var userID string
		if dataMap["user_id"] != nil {
			userID = dataMap["user_id"].(string)
		} else {
			userID = ""
		}

		if len(userID) == 0 {
			fmt.Println("storeRequestLimitMiddleware: empty user_id -> go next")
			return next(data, metadata)
		}

		checkReq := MiddlewareCheckRequest{
			Method: "check",
			Data: MiddlewareCheckRequestData{
				RequestMicroservice: microserviceName,
				RequestMethod:       method,
				StoreUserID:         userID,
			},
		}

		jsonData, err := json.Marshal(checkReq)
		if err != nil {
			fmt.Println("storeRequestLimitMiddleware: error marshaling data")
			fmt.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("marshaling -> " + err.Error())
		}

		req, err := http.NewRequest("POST", requestLimitServiceURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("storeRequestLimitMiddleware: error creating request")
			fmt.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("creating request -> " + err.Error())
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("storeRequestLimitMiddleware: error sending request to auth")
			fmt.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("sending request -> " + err.Error())
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("storeRequestLimitMiddleware: error reading body from auth")
			fmt.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("reading body -> " + err.Error())
		}

		var response CheckResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("storeRequestLimitMiddleware: error unmarshalling body from auth")
			fmt.Println("storeRequestLimitMiddleware: " + err.Error())
			return unauthorizedResponse("Unmarshal -> " + err.Error())
		}

		if response.Status != ResponseStatusOK {
			fmt.Println("storeRequestLimitMiddleware: response-body -> result is not `Ok`")
			fmt.Println("storeRequestLimitMiddleware: " + string(body))
			return unauthorizedResponse("Result -> " + string(body))
		} else if !response.Result.IsRequestAvailable {
			return unauthorizedResponse(" request limit exceeded")
		}

		return next(data, metadata)
	}
}
