package golrackpi

import (
	//"bytes"
	"encoding/json"
	"errors"

	"io/ioutil"
	"net/http"
	//"sort"
)

func (c *AuthClient) Version() (map[string]interface{}, error) {
	var result map[string]interface{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/info/version"), nil)
	if err != nil {
		return result, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return result, errMe

	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	var jsonResult interface{}
	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		//fmt.Println(errJson)
		return result, errJson
	}
	//fmt.Println(jsonResult)

	m, mOk := jsonResult.(map[string]interface{})

	if mOk {
		return m, nil
		// Use Map
	}
	return result, errors.New("could not read response")

}
