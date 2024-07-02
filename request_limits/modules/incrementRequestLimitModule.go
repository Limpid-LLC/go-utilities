package request_limits_modules

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
)

var IncrementRequestLimitModule *IncrementRequestLimitModuleObj

type IncrementRequestLimitModuleObj struct {
	RequestLimitServiceURL string
	MicroserviceName       string
	MasterToken            string
}

func InitIncrementRequestLimitModule(requestLimitServiceURL string, microserviceName string, masterToken string) {
	IncrementRequestLimitModule = &IncrementRequestLimitModuleObj{
		RequestLimitServiceURL: requestLimitServiceURL,
		MicroserviceName:       microserviceName,
		MasterToken:            masterToken,
	}
}

func (module *IncrementRequestLimitModuleObj) Increment(
	microserviceMethod string,
	serviceStationID string,
	storeUserID string,
	date int64,
) error {
	if module.RequestLimitServiceURL == "" {
		err := errors.New("IncrementRequestLimitModule: requestLimitServiceURL is empty")
		fmt.Println(err.Error())
		return err
	}
	if module.MicroserviceName == "" {
		err := errors.New("IncrementRequestLimitModule: microserviceName is empty")
		fmt.Println(err.Error())
		return err
	}

	if len(serviceStationID) == 0 && len(storeUserID) == 0 {
		err := errors.New("IncrementRequestLimitModule: serviceStationID or storeUserID required")
		fmt.Println(err.Error())
		return err
	}

	if date == 0 {
		err := errors.New("IncrementRequestLimitModule: date is required")
		fmt.Println(err.Error())
		return err
	}

	decrementReq := MicroserviceIncrementRequest{
		Method: "increment_request_limit",
		Metadata: bson.M{
			"token": module.MasterToken,
		},
		Data: IncrementRequestLimitRequestData{
			RequestMicroservice: module.MicroserviceName,
			RequestMethod:       microserviceMethod,
			ServiceStationID:    serviceStationID,
			StoreUserID:         storeUserID,
			Date:                date,
		},
	}

	jsonData, err := json.Marshal(decrementReq)
	if err != nil {
		fmt.Println("IncrementRequestLimitModule: error marshaling data")
		fmt.Println("IncrementRequestLimitModule: " + err.Error())
		return err
	}

	req, err := http.NewRequest("POST", module.RequestLimitServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("IncrementRequestLimitModule: error creating request")
		fmt.Println("IncrementRequestLimitModule: " + err.Error())
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("IncrementRequestLimitModule: error sending request to request limit service")
		fmt.Println("IncrementRequestLimitModule: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("IncrementRequestLimitModule: error reading body from request limit service")
		fmt.Println("IncrementRequestLimitModule: " + err.Error())
		return err
	}

	var response MicroserviceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("IncrementRequestLimitModule: error unmarshalling body from request limit service")
		fmt.Println("IncrementRequestLimitModule: " + err.Error())
		return err
	}

	if response.Status != ResponseStatusOK {
		fmt.Println("IncrementRequestLimitModule: response-body -> result is not `Ok`")
		fmt.Println("IncrementRequestLimitModule: " + string(body))
		return errors.New("IncrementRequestLimitModule Response Error: " + response.Error)
	}

	return nil
}
