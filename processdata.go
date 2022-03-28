package golrackpi

import (
	"bytes"
	"encoding/json"

	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
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

func (c *AuthClient) GetProcessDataList() (map[string]ProcessData, error) {
	processData := make(map[string]ProcessData)
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
	//sb := string(body)
	//fmt.Println("raw body output:")
	//fmt.Println(sb)

	//fmt.Println(response.Body)
	//var resultMe map[string]interface{}
	var jsonResult interface{}
	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		//fmt.Println(errJson)
		return processData, errJson
	}
	//fmt.Println(jsonResult)

	s, sOk := jsonResult.([]interface{})
	if !sOk {

		return processData, errors.New("error in conversion of JSON result")
	}

	// Use Slice
	//fmt.Println("use slice")
	//fmt.Println(s)
	for _, v := range s {
		//fmt.Println(v)
		switch vv := v.(type) {
		case map[string]interface{}:

			moduleid := vv["moduleid"].(string)
			processdataids := vv["processdataids"].([]interface{})
			//fmt.Println("moduleid:", moduleid)
			//fmt.Println("processdataids", processdataids)

			var processDataIds []string
			for _, p := range processdataids {
				//fmt.Println("i, p:", i, p)
				processDataIds = append(processDataIds, p.(string))

			}
			sort.Strings(processDataIds)
			processData[moduleid] = ProcessData{ModuleId: moduleid, ProcessDataIds: processDataIds}

			//fmt.Println(vv["moduleid"])
			//c.writeJson(vv)
		default:

			return processData, errors.New("unknown returned type in inverter response")
		}
	}

	//fmt.Println("Result:", processData)

	return processData, nil

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

				moduleid := vv["moduleid"].(string)
				processdata := vv["processdata"].([]interface{})
				fmt.Println("moduleid:", moduleid)
				fmt.Println("processdata", processdata)
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

}

func (c *AuthClient) GetProcessDataValues(v []ProcessData) map[string]ProcessDataValues {

	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(string(b))

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
	//sb := string(body)
	//fmt.Println("raw body output:")
	//fmt.Println(sb)

	var jsonResult interface{}
	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		fmt.Println(errJson)
	}
	fmt.Println(jsonResult)

	m, mOk := jsonResult.(map[string]interface{})
	s, _ := jsonResult.([]interface{})

	//m := jsonResult.(map[string]interface{})
	var processDataValues ProcessDataValues
	resultData := make(map[string]ProcessDataValues)

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
}