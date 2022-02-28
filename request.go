package golrackpi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *AuthClient) Modules() {
	client := http.Client{}

	//request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/modules"), nil)
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
				c.writeJson(vv)
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

}

func (c *AuthClient) ProcessData() {
	client := http.Client{}

	//request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata/devices:local/HomeOwn_P"), nil)
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
				c.writeJson(vv)
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

	/*for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}*/
	//json.NewDecoder(response.Body).Decode(&resultMe)
	//fmt.Println(resultMe)

}

func (c *AuthClient) writeJson(data map[string]interface{}) {
	fmt.Println("in writeJson")
	for k, v := range data {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}

func (c *AuthClient) Settings() {
	client := http.Client{}

	//request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/settings"), nil)
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
				c.writeJson(vv)
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

}
