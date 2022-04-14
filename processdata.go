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

// ProcessData specifies the structure of the response returned by a request to the "processdata" endpoint
type ProcessData struct {
	ModuleId       string   `json:"moduleid"`
	ProcessDataIds []string `json:"processdataids"`
}

// ProcessDataValue specifies the structure of a single processdata value with its Unit, Processdata-ID and Value
type ProcessDataValue struct {
	Unit  string      `json:"unit"`
	Id    string      `json:"id"`
	Value interface{} `json:"value"`
}

// ProcessDataValues specifies the structure of the response returned by a request for processdata with moduleid and
// a slice of ProcessDataValue which contains the fields Unid, Id and Value
type ProcessDataValues struct {
	ModuleId    string             `json:"moduleid"`
	ProcessData []ProcessDataValue `json:"processdata"`
}

// ProcessData returns a slice of ProcessData type, i.e. a list of modules with a list of their process-data identifiers
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
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return processData, err
	}

	errJson := json.Unmarshal(body, &processData)
	if errJson != nil {
		return processData, errJson
	}

	return processData, nil

}

// ProcessDataModule returns a slice of ProcessDataValues returned by the request to the "processdata/moduleid" endpoint.
// It takes a moduleid and returns all processdata ids and their values according to the moduleid.
func (c *AuthClient) ProcessDataModule(moduleId string) ([]ProcessDataValues, error) {
	processDataValues := []ProcessDataValues{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata/"+moduleId), nil)
	if err != nil {
		return processDataValues, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return processDataValues, errMe
	}
	defer response.Body.Close()

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

// ProcessDataModuleValues returns a slice of ProcessDataValues returned by a request of moduleid and one or more of the processdataids which
// is handled by the module.
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
	defer response.Body.Close()

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
	defer response.Body.Close()

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
