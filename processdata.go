package golrackpi

import (
	"bytes"
	"encoding/json"
	"strings"

	"io/ioutil"
	"net/http"
	//"time"
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

func (c *AuthClient) ProcessData() ([]ProcessData, error) {
	processData := []ProcessData{}
	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	if err != nil {
		return processData, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {

		return processData, errMe
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return processData, err
	}

	//fmt.Println(response.Body)
	//var resultMe map[string]interface{}
	errJson := json.Unmarshal(body, &processData)
	if errJson != nil {
		//fmt.Println(errJson)
		return processData, errJson
	}

	return processData, nil

}

func (c *AuthClient) ProcessDataModuleValues(moduleId string, processDataIds ...string) ([]ProcessDataValues, error) {
	processDataValues := []ProcessDataValues{}
	client := http.Client{}

	var processDataString string

	if len(processDataIds) > 1 { // ok, this is not really necessary...
		processDataString = strings.TrimRight(strings.Join(processDataIds, ","), ",")

	} else {
		processDataString = processDataIds[0]
	}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata/"+moduleId+"/"+processDataString), nil)
	if err != nil {
		return processDataValues, err

	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return processDataValues, errMe
	}

	body, err := ioutil.ReadAll(response.Body)
	//sb := string(body)
	//fmt.Println("raw body output:")
	//fmt.Println(sb)

	errJson := json.Unmarshal(body, &processDataValues)
	if errJson != nil {
		return processDataValues, errJson

	}

	return processDataValues, nil

}

func (c *AuthClient) ProcessDataValues(v []ProcessData) ([]ProcessDataValues, error) {
	processDataValues := []ProcessDataValues{}
	b, err := json.Marshal(v)
	if err != nil {
		return processDataValues, err

	}

	//fmt.Println(string(b))

	client := http.Client{}

	request, err := http.NewRequest("POST", c.getUrl("/api/v1/processdata"), bytes.NewBuffer(b))
	if err != nil {
		return processDataValues, err

	}
	request.Header.Add("Content-Type", "application/json")

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errReq := client.Do(request)
	if errReq != nil {
		return processDataValues, errReq
	}
	//fmt.Println(response)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return processDataValues, err

	}
	//sb := string(body)
	//fmt.Println("raw body output:")
	//fmt.Println(sb)

	errJson := json.Unmarshal(body, &processDataValues)
	if errJson != nil {

		return processDataValues, errJson
	}
	//fmt.Println(processDataValues)

	return processDataValues, nil
	//m := jsonResult.(map[string]interface{})
	/*

		var moduleid string
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

					moduleid = vv["moduleid"].(string)
					processdata := vv["processdata"].([]interface{})
					fmt.Println("moduleid:", moduleid)
					fmt.Println("processdata", processdata)

					var processDataValue []ProcessDataValue
					//var processDataIds []string
					for i, p := range processdata {
						fmt.Println("i, p:", i, p)

						d := p.(map[string]interface{})
						fmt.Println("data", d)

						fmt.Println("Unit:", d["unit"])
						fmt.Println("Id:", d["id"])

						fmt.Println("Value:", d["value"])

						pdValue := ProcessDataValue{Unit: d["unit"].(string), Id: d["id"].(string), Value: d["value"]}
						fmt.Println("pdValue", pdValue)
						processDataValue = append(processDataValue, pdValue)
						//processDataIds = append(processDataIds, p.(string))
						//processData[moduleid].ProcessDataIds = append(processData[moduleid].ProcessDataIds, p)
					}
					//sort.Strings(processDataIds)
					//processData[moduleid] = ProcessData{ModuleId: moduleid, ProcessDataIds: processDataIds}
					processDataValues = ProcessDataValues{ModuleId: moduleid, ProcessData: processDataValue}
					fmt.Println("result", processDataValues)

					//c.writeJson(vv)
				case []interface{}:
					fmt.Println(k, "is an array:")
					for i, u := range vv {
						fmt.Println(i, u)
					}
				default:
					fmt.Println(k, "is of a type I don't know how to handle", vv)
				}
				resultData[moduleid] = processDataValues
			}

		}
		return resultData
	*/
}
