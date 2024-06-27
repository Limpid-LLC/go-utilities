package request_limits_modules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
)

var DecrementRequestLimitModule *DecrementRequestLimitModuleObj

type DecrementRequestLimitModuleObj struct {
	RequestLimitServiceURL string
	MicroserviceName       string
	MasterToken            string
}

func InitDecrementRequestLimitModule(requestLimitServiceURL string, microserviceName string, masterToken string) {
	DecrementRequestLimitModule = &DecrementRequestLimitModuleObj{
		RequestLimitServiceURL: requestLimitServiceURL,
		MicroserviceName:       microserviceName,
		MasterToken:            masterToken,
	}
}

func (module *DecrementRequestLimitModuleObj) Decrement(
	microserviceMethod string,
	serviceStationID string,
	storeUserID string,
) error {
	if module.RequestLimitServiceURL == "" {
		err := errors.New("DecrementRequestLimitModule: requestLimitServiceURL is empty")
		fmt.Println(err.Error())
		return err
	}
	if module.MicroserviceName == "" {
		err := errors.New("DecrementRequestLimitModule: microserviceName is empty")
		fmt.Println(err.Error())
		return err
	}

	if len(serviceStationID) == 0 && len(storeUserID) == 0 {
		err := errors.New("DecrementRequestLimitModule: serviceStationID or storeUserID required")
		fmt.Println(err.Error())
		return err
	}

	decrementReq := MicroserviceRequest{
		Method: "decrement_request_limit",
		Metadata: bson.M{
			"token": module.MasterToken,
		},
		Data: DecrementRequestLimitRequestData{
			RequestMicroservice: module.MicroserviceName,
			RequestMethod:       microserviceMethod,
			ServiceStationID:    serviceStationID,
			StoreUserID:         storeUserID,
		},
	}

	jsonData, err := json.Marshal(decrementReq)
	if err != nil {
		fmt.Println("DecrementRequestLimitModule: error marshaling data")
		fmt.Println("DecrementRequestLimitModule: " + err.Error())
		return err
	}

	req, err := http.NewRequest("POST", module.RequestLimitServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("DecrementRequestLimitModule: error creating request")
		fmt.Println("DecrementRequestLimitModule: " + err.Error())
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("DecrementRequestLimitModule: error sending request to request limit service")
		fmt.Println("DecrementRequestLimitModule: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("DecrementRequestLimitModule: error reading body from  request limit service")
		fmt.Println("DecrementRequestLimitModule: " + err.Error())
		return err
	}

	var response MicroserviceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("DecrementRequestLimitModule: error unmarshalling body from  request limit service")
		fmt.Println("DecrementRequestLimitModule: " + err.Error())
		return err
	}

	if response.Status != ResponseStatusOK {
		fmt.Println("DecrementRequestLimitModule: response-body -> result is not `Ok`")
		fmt.Println("DecrementRequestLimitModule: " + string(body))
		return errors.New("DecrementRequestLimitModule Response Error: " + response.Error)
	}

	return nil
}
