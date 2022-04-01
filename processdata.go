// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golrackpi

import (
	"bytes"
	"encoding/json"
	"strings"

	"io/ioutil"
	"net/http"
)

type ProcessData struct {
	ModuleId       string   `json:"moduleid"`
	ProcessDataIds []string `json:"processdataids"`
}

type ProcessDataValue struct {
	Unit  string      `json:"unit"`
	Id    string      `json:"id"`
	Value interface{} `json:"value"`
}

type ProcessDataValues struct {
	ModuleId    string             `json:"moduleid"`
	ProcessData []ProcessDataValue `json:"processdata"`
}

func (c *AuthClient) ProcessData() ([]ProcessData, error) {
	processData := []ProcessData{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	if err != nil {
		return processData, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {

		return processData, errMe
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return processData, err
	}

	//fmt.Println(response.Body)
	//var resultMe map[string]interface{}
	errJson := json.Unmarshal(body, &processData)
	if errJson != nil {
		//fmt.Println(errJson)
		return processData, errJson
	}

	return processData, nil

}

func (c *AuthClient) ProcessDataModuleValues(moduleId string, processDataIds ...string) ([]ProcessDataValues, error) {
	processDataValues := []ProcessDataValues{}
	client := http.Client{}

	var processDataString string

	if len(processDataIds) > 1 {
		processDataString = strings.TrimRight(strings.Join(processDataIds, ","), ",")

	} else { // ok, this is not really necessary...
		processDataString = processDataIds[0]
	}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata/"+moduleId+"/"+processDataString), nil)
	if err != nil {
		return processDataValues, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return processDataValues, errMe
	}

	body, err := ioutil.ReadAll(response.Body)

	errJson := json.Unmarshal(body, &processDataValues)
	if errJson != nil {
		return processDataValues, errJson
	}

	return processDataValues, nil

}

func (c *AuthClient) ProcessDataValues(v []ProcessData) ([]ProcessDataValues, error) {
	processDataValues := []ProcessDataValues{}
	b, err := json.Marshal(v)
	if err != nil {
		return processDataValues, err
	}

	client := http.Client{}

	request, err := http.NewRequest("POST", c.getUrl("/api/v1/processdata"), bytes.NewBuffer(b))
	if err != nil {
		return processDataValues, err
	}
	request.Header.Add("Content-Type", "application/json")

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errReq := client.Do(request)
	if errReq != nil {
		return processDataValues, errReq
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return processDataValues, err
	}

	errJson := json.Unmarshal(body, &processDataValues)
	if errJson != nil {
		return processDataValues, errJson
	}

	return processDataValues, nil

}
