// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golrackpi

import (
	"bytes"
	"encoding/json"
	"io"

	"net/http"

	"github.com/geschke/golrackpi/internal/timefix"
)

// EventData specifies the structure of the response returned by a request to the "events" endpoint
type EventData struct {
	Description     string               `json:"description"`
	Category        string               `json:"category"`
	LongDescription string               `json:"long_description"`
	StartTime       timefix.InverterTime `json:"start_time"`
	Group           string               `json:"group"`
	EndTime         timefix.InverterTime `json:"end_time"`
	Code            int                  `json:"code"`
	IsActive        bool                 `json:"is_active"`
}

// EventsWithParam returns the latest events with localized descriptions. It returns a slice of EventData type. It takes as arguments
// the language string (currently available de-de, en-gb, es-es, fr-fr, hu-hu, it-it, nl-nl, pl-pl, pt-pt, cs-cz, el-gr and zh-cn) and
// the maximum number of events (default: 10)
func (c *AuthClient) EventsWithParam(language string, max int) ([]EventData, error) {
	jsonResult := []EventData{}
	if language == "" {
		language = "en-gb"
	}
	if max <= 0 {
		max = 10
	}

	payload := struct {
		Language string `json:"language"`
		Max      int    `json:"max"`
	}{
		Language: language,
		Max:      max,
	}

	b, err := json.Marshal(payload)
	if err != nil {

		return jsonResult, err
	}

	client := http.Client{}

	request, err := http.NewRequest("POST", c.getUrl("/api/v1/events/latest"), bytes.NewBuffer(b))
	if err != nil {
		return jsonResult, err
	}
	request.Header.Add("Content-Type", "application/json")

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

// Events returns the latest events as a slice of EventData type
func (c *AuthClient) Events() ([]EventData, error) {
	jsonResult := []EventData{}

	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/events/latest"), nil)
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
