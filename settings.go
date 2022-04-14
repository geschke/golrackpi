// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golrackpi

import (
	"encoding/json"
	"strings"

	"errors"
	"io/ioutil"
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

	response, errMe := client.Do(request)
	if errMe != nil {
		return jsonResult, errMe
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		return jsonResult, errJson
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

	response, errMe := client.Do(request)
	if errMe != nil {
		return jsonResult, errMe
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		return jsonResult, errJson

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

	response, errMe := client.Do(request)
	if errMe != nil {
		return jsonResult, errMe
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		return jsonResult, errJson

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

	response, errMe := client.Do(request)
	if errMe != nil {
		return jsonResult, errMe
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		return jsonResult, errJson

	}

	return jsonResult, nil
}
