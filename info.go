package golrackpi

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"sort"
)

func (c *AuthClient) Version() map[string]interface{} {
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/info/version"), nil)
	if err != nil {
		fmt.Println(err)
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		fmt.Println(errMe)
	}

	body, err := ioutil.ReadAll(response.Body)
	sb := string(body)
	fmt.Println("raw body output:")
	fmt.Println(sb)

	fmt.Println(response.Body)
	//var resultMe map[string]interface{}
	var jsonResult interface{}
	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		fmt.Println(errJson)
	}
	fmt.Println(jsonResult)

	m, mOk := jsonResult.(map[string]interface{})
	//s, _ := jsonResult.([]interface{})

	//moduleData := make(map[string]ModuleData)

	//m := jsonResult.(map[string]interface{})
	if mOk {
		return m
		// Use Map
	} else {
		// error
	}
	return m
	//return moduleData
}
