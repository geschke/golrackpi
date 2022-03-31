package golrackpi

import (
	"encoding/json"

	"io/ioutil"
	"net/http"
	//"time"
)

type ModuleData struct {
	Id   string
	Type string
}

func (c *AuthClient) Modules() ([]ModuleData, error) {
	moduleData := []ModuleData{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/modules"), nil)
	if err != nil {
		return moduleData, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return moduleData, errMe
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return moduleData, err

	}
	errJson := json.Unmarshal(body, &moduleData)
	if errJson != nil {
		return moduleData, errJson

	}

	return moduleData, nil
}
