// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golrackpi

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"net/http"
)

// ProcessData specifies the structure of the response returned by a request to the "processdata" endpoint.
// Furthermore this type is used in the ProcessDataValues function to define an arbitrary number of moduleids and
// also an arbitrary number of their processdataids.
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

	response, err := client.Do(request)
	if err != nil {
		return processData, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return processData, err
	}

	err = json.Unmarshal(body, &processData)
	if err != nil {
		return processData, err
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

	response, err := client.Do(request)
	if err != nil {
		return processDataValues, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return processDataValues, err
	}

	err = json.Unmarshal(body, &processDataValues)
	if err != nil {
		return processDataValues, err
	}

	return processDataValues, nil

}

// ProcessDataModuleValues returns a slice of ProcessDataValues returned by a request of a moduleid and one or more of the processdataids which
// belongs to the moduleid.
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

	response, err := client.Do(request)
	if err != nil {
		return processDataValues, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return processDataValues, err
	}

	err = json.Unmarshal(body, &processDataValues)
	if err != nil {
		return processDataValues, err
	}

	return processDataValues, nil

}

// ProcessDataValues returns a slice of ProcessDataValues by a request of an arbitrary number of modules with one or more processdataids
// according to the moduleid.
// It takes a slice of ProcessData as argument, so it's possible to submit several moduleids with an arbitrary number of their processdataids
// and get all processdata values with one request to the inverter.
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

	response, err := client.Do(request)
	if err != nil {
		return processDataValues, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return processDataValues, err
	}

	err = json.Unmarshal(body, &processDataValues)
	if err != nil {
		return processDataValues, err
	}

	return processDataValues, nil

}
