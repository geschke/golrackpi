// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golrackpi

import (
	"encoding/json"
	"errors"
	"io"

	"net/http"
)

// Version returns information about the API; currently name, hostname, sw_version and api_version
func (c *AuthClient) Version() (map[string]interface{}, error) {
	var result map[string]interface{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/info/version"), nil)
	if err != nil {
		return result, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, err := client.Do(request)
	if err != nil {
		return result, err

	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	var jsonResult interface{}
	err = json.Unmarshal(body, &jsonResult)
	if err != nil {
		return result, err
	}

	m, mOk := jsonResult.(map[string]interface{})

	if mOk {
		return m, nil
	}
	return result, errors.New("could not read response")

}
