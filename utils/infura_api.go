package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

var ApiURL = os.Getenv("API_BASE_URL") + os.Getenv("INFURA_PROJECT_ID")

func GetAPIBodyRequest(method string, params []string) (*bytes.Buffer, error) {
	body := map[string]interface{}{
		"jsonrpc": "2.0", "method": method, "params": params, "id": 1,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonData), nil
}

func InfuraAPIRequest(body *bytes.Buffer) (map[string]interface{}, error) {
	response, err := http.Post(ApiURL, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(responseBody, &data)
	return data, err
}
