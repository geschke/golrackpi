package golrackpi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type ProcessData struct {
	ModuleId       string   `json:"moduleid"`
	ProcessDataIds []string `json:"processdataids"`
}

type ProcessDataValue struct {
	Unit  string      `json:"unit"`
	Id    string      `json:"id"`
	Value interface{} `json:"value"`
}

type ProcessDataValues struct {
	ModuleId    string             `json:"moduleid"`
	ProcessData []ProcessDataValue `json:"processdata"`
}

type ModuleData struct {
	Id   string
	Type string
}

func (c *AuthClient) Modules() map[string]ModuleData {
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

	moduleData := make(map[string]ModuleData)

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
				fmt.Println("moduleid:", moduleId)
				fmt.Println("typeData", typeData)

				/*var processDataIds []string
				for i, p := range processdataids {
					fmt.Println("i, p:", i, p)
					processDataIds = append(processDataIds, p.(string))
					//processData[moduleid].ProcessDataIds = append(processData[moduleid].ProcessDataIds, p)
				}
				sort.Strings(processDataIds)*/
				moduleData[moduleId] = ModuleData{Id: moduleId, Type: typeData}

				//c.writeJson(vv)
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
	return moduleData
}

func (c *AuthClient) GetProcessDataList() map[string]ProcessData {
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	//request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata/devices:local/HomeOwn_P"), nil)
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

	processData := make(map[string]ProcessData)

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
				moduleid := vv["moduleid"].(string)
				processdataids := vv["processdataids"].([]interface{})
				fmt.Println("moduleid:", moduleid)
				fmt.Println("processdataids", processdataids)

				var processDataIds []string
				for i, p := range processdataids {
					fmt.Println("i, p:", i, p)
					processDataIds = append(processDataIds, p.(string))
					//processData[moduleid].ProcessDataIds = append(processData[moduleid].ProcessDataIds, p)
				}
				sort.Strings(processDataIds)
				processData[moduleid] = ProcessData{ModuleId: moduleid, ProcessDataIds: processDataIds}

				fmt.Println(vv["moduleid"])
				//c.writeJson(vv)
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

	fmt.Println("Result:", processData)

	return processData
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

func (c *AuthClient) GetProcessData(moduleId string, processDataId string) {
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata/"+moduleId+"/"+processDataId), nil)
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

func (c *AuthClient) GetProcessDataValues() {
	// try to build appropriate JSON

	pd := []string{"Statistic:Autarky:Day",
		"Statistic:Autarky:Month",
		"Statistic:Autarky:Total"}
	v := ProcessData{ModuleId: "scb:statistic:EnergyFlow", ProcessDataIds: pd}

	pd2 := []string{"HomeBat_P", "HomeGrid_P"}
	v2 := ProcessData{ModuleId: "devices:local", ProcessDataIds: pd2}

	v3 := []ProcessData{v, v2}

	b, err := json.Marshal(v3)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(string(b))
	/*[
	  {
	    "processdataids": [
	      "Statistic:Autarky:Day",
	      "Statistic:Autarky:Month",
	      "Statistic:Autarky:Total"
	    ],
	    "moduleid": "scb:statistic:EnergyFlow"
	  },
	  {
	    "processdataids": [
	    "HomeBat_P",
	    "HomeGrid_P"
	    ],
	    "moduleid": "devices:local" }
	]*/
	client := http.Client{}

	request, err := http.NewRequest("POST", c.getUrl("/api/v1/processdata"), bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Add("Content-Type", "application/json")

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errReq := client.Do(request)
	if errReq != nil {
		fmt.Println(errReq)
	}
	//fmt.Println(response)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	sb := string(body)
	fmt.Println("raw body output:")
	fmt.Println(sb)

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