// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golrackpi

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"errors"

	"net/http"
)

// SettingsDataValues specifies the structure of a setting element returned by a request to the "settings" endpoint
type SettingsDataValues struct {
	Id      string      `json:"id"`
	Max     string      `json:"max"`
	Min     string      `json:"min"`
	Unit    string      `json:"unit"`
	Type    interface{} `json:"type"`
	Access  interface{} `json:"access"`
	Default string      `json:"default"`
}

// SettingsData specifies the structure of the response returned by a request to the "settings" endpoint.
// It embeds a slice of SettingsDataValues type.
type SettingsData struct {
	ModuleId string               `json:"moduleid"`
	Settings []SettingsDataValues `json:"settings"`
}

// SettingsValues specifies the structure of the response returned by a request to the "settings/moduleid/..." endpoint.
// The structure defines a settingid and its value.
type SettingsValues struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

// Settings returns a list of all modules with their setting identifiers and further parameters of the setting, i.e. max, min, default etc.
// Warning: The request returns a lot of data, so it takes some time.
func (c *AuthClient) Settings() ([]SettingsData, error) {
	jsonResult := []SettingsData{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/settings"), nil)
	if err != nil {
		return jsonResult, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, err := client.Do(request)
	if err != nil {
		return jsonResult, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	err = json.Unmarshal(body, &jsonResult)
	if err != nil {
		return jsonResult, err
	}

	return jsonResult, nil
}

// SettingsModule returns a list of settings with settingids and their values of a moduleid
func (c *AuthClient) SettingsModule(moduleid string) ([]SettingsValues, error) {
	jsonResult := []SettingsValues{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/settings/"+moduleid), nil)
	if err != nil {
		return jsonResult, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, err := client.Do(request)
	if err != nil {
		return jsonResult, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	err = json.Unmarshal(body, &jsonResult)
	if err != nil {
		return jsonResult, err

	}

	return jsonResult, nil
}

// SettingsModuleSetting returns a SettingsValues slice with length 1 according to the submitted
// moduleid and settingid parameter.
func (c *AuthClient) SettingsModuleSetting(moduleid string, settingid string) ([]SettingsValues, error) {
	jsonResult := []SettingsValues{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/settings/"+moduleid+"/"+settingid), nil)
	if err != nil {
		return jsonResult, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, err := client.Do(request)
	if err != nil {
		return jsonResult, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	err = json.Unmarshal(body, &jsonResult)
	if err != nil {
		return jsonResult, err

	}

	return jsonResult, nil
}

// SettingsModuleSettings returns a SettingsValues slice according to the submitted
// moduleid and settingids parameter. This function takes an arbitrary number of setting ids as arguments.
func (c *AuthClient) SettingsModuleSettings(moduleid string, settingids ...string) ([]SettingsValues, error) {
	jsonResult := []SettingsValues{}
	client := http.Client{}

	for i, settingid := range settingids {
		settingids[i] = strings.TrimSpace(settingid)
	}
	csvSettings := strings.Join(settingids, ",")

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/settings/"+moduleid+"/"+csvSettings), nil)
	if err != nil {
		return jsonResult, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, err := client.Do(request)
	if err != nil {
		return jsonResult, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	err = json.Unmarshal(body, &jsonResult)
	if err != nil {
		return jsonResult, err

	}

	return jsonResult, nil
}

// ModuleSettings specifies the structure of the request body for the "settings" endpoint.
type ModuleSettings struct {
	Settings []SettingsValues `json:"settings"`
	ModuleId string           `json:"moduleid"`
}

// write settings
func (c *AuthClient) UpdateSettings(settings []ModuleSettings) ([]ModuleSettings, error) {
	jsonResult := []ModuleSettings{}
	jsonPayload, err := json.Marshal(settings)

	if err != nil {
		return jsonResult, err
	}

	client := http.Client{}
	request, err := http.NewRequest("PUT", c.getUrl("/api/v1/settings"), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return jsonResult, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return jsonResult, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return jsonResult, err
	}

	body, err := io.ReadAll(response.Body) // response body is []byte
	if err != nil {
		return jsonResult, err
	}

	err = json.Unmarshal(body, &jsonResult)
	if err != nil {
		return jsonResult, err

	}

	return jsonResult, err
}
