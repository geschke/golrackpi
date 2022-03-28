package golrackpi

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"
	//"time"
)

type ModuleData struct {
	Id   string
	Type string
}

func (c *AuthClient) Modules() (map[string]ModuleData, error) {
	moduleData := make(map[string]ModuleData)
	client := http.Client{}

	//request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

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
	//sb := string(body)
	//fmt.Println("raw body output:")
	//fmt.Println(sb)

	//fmt.Println(response.Body)
	//var resultMe map[string]interface{}
	var jsonResult interface{}
	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		return moduleData, errJson

	}
	//fmt.Println(jsonResult)

	m, mOk := jsonResult.(map[string]interface{})
	s, _ := jsonResult.([]interface{})

	//m := jsonResult.(map[string]interface{})
	if mOk {
		// Use Map
		fmt.Println("use map")
		fmt.Println(m)
	} else {
		// Use Slice
		fmt.Println("use slice")
		fmt.Println(s)
		for k, v := range s {
			fmt.Println(k)
			fmt.Println(v)
			switch vv := v.(type) {
			case string:
				fmt.Println(k, "is string", vv)
			case float64:
				fmt.Println(k, "is float64", vv)
			case map[string]interface{}:
				fmt.Println(k, "is map dingens", vv)

				moduleId := vv["id"].(string)
				typeData := vv["type"].(string)

				moduleData[moduleId] = ModuleData{Id: moduleId, Type: typeData}

			case []interface{}:
				fmt.Println(k, "is an array:")
				for i, u := range vv {
					fmt.Println(i, u)
				}
			default:
				fmt.Println(k, "is of a type I don't know how to handle", vv)
			}
		}
	}
	return moduleData, nil
}
