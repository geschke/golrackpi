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

type SettingsDataValues struct {
	Id      string      `json:"id"`
	Max     string      `json:"max"`
	Min     string      `json:"min"`
	Unit    string      `json:"unit"`
	Type    interface{} `json:"type"`
	Access  interface{} `json:"access"`
	Default string      `json:"default"`
}

type SettingsData struct {
	ModuleId string               `json:"moduleid"`
	Settings []SettingsDataValues `json:"settings"`
}

type SettingsValues struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

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
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
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
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
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

func (c *AuthClient) SettingsModuleSettings(moduleid string, settingids []string) ([]SettingsValues, error) {
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
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
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
