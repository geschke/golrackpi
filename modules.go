// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golrackpi

import (
	"encoding/json"
	"io"

	"net/http"
)

// ModuleData specifies the structure of the response returned by a request to the "modules" endpoint
type ModuleData struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

// Modules returns a list of modules with their type
func (c *AuthClient) Modules() ([]ModuleData, error) {
	moduleData := []ModuleData{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/modules"), nil)
	if err != nil {
		return moduleData, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, err := client.Do(request)
	if err != nil {
		return moduleData, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return moduleData, err

	}
	err = json.Unmarshal(body, &moduleData)
	if err != nil {
		return moduleData, err

	}

	return moduleData, nil
}
