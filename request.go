package golrackpi

import (
	"bytes"
	"encoding/json"
	"strings"

	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	//"time"
	"github.com/geschke/golrackpi/internal/timefix"
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

type SettingsDataValues struct {
	Id      string      `json:"id"`
	Max     string      `json:"max"`
	Min     string      `json:"min"`
	Unit    string      `json:"unit"`
	Type    interface{} `json:"type"`
	Access  interface{} `json:"access"`
	Default string      `json:"default"`
}

type SettingsData struct {
	ModuleId string               `json:"moduleid"`
	Settings []SettingsDataValues `json:"settings"`
}

type SettingsValues struct {
	Id    string `json:"id"`
	Value string `json:"value"`
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

func (c *AuthClient) Settings() ([]SettingsData, error) {
	jsonResult := []SettingsData{}
	client := http.Client{}

	//request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/settings"), nil)
	if err != nil {
		return jsonResult, err

	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return jsonResult, errMe

	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}
	//sb := string(body)
	//fmt.Println("raw body output:")
	//fmt.Println(sb)

	//fmt.Println(response.Body)
	//var resultMe map[string]interface{}

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		return jsonResult, errJson

	}
	//fmt.Println(jsonResult)
	return jsonResult, nil
	/*	m, mOk := jsonResult.(map[string]interface{})
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
	*/
}

func (c *AuthClient) SettingsModule(moduleid string) ([]SettingsValues, error) {
	jsonResult := []SettingsValues{}
	client := http.Client{}

	//request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/settings/"+moduleid), nil)
	if err != nil {
		return jsonResult, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return jsonResult, errMe
	}
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	sb := string(body)
	fmt.Println("raw body output:")
	fmt.Println(sb)

	//fmt.Println(response.Body)

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		return jsonResult, errJson

	}

	return jsonResult, nil
}

func (c *AuthClient) SettingsModuleSetting(moduleid string, settingid string) ([]SettingsValues, error) {
	jsonResult := []SettingsValues{}
	client := http.Client{}

	//request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/settings/"+moduleid+"/"+settingid), nil)
	if err != nil {
		return jsonResult, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return jsonResult, errMe
	}
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	sb := string(body)
	fmt.Println("raw body output:")
	fmt.Println(sb)

	//fmt.Println(response.Body)

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		return jsonResult, errJson

	}

	return jsonResult, nil
}

func (c *AuthClient) SettingsModuleSettings(moduleid string, settingids []string) ([]SettingsValues, error) {
	jsonResult := []SettingsValues{}
	client := http.Client{}

	//request, err := http.NewRequest("GET", c.getUrl("/api/v1/processdata"), nil)

	for i, settingid := range settingids {
		settingids[i] = strings.TrimSpace(settingid)
	}
	csvSettings := strings.Join(settingids, ",")

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/settings/"+moduleid+"/"+csvSettings), nil)
	if err != nil {
		return jsonResult, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return jsonResult, errMe
	}
	if response.StatusCode != 200 {
		return jsonResult, errors.New("module or setting not found")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err
	}

	sb := string(body)
	fmt.Println("raw body output:")
	fmt.Println(sb)

	//fmt.Println(response.Body)

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		return jsonResult, errJson

	}

	return jsonResult, nil
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

// EventsCustomized returns the latest events with localized descriptions. It takes as arguments
// the language string (currently available de-de, en-gb, es-es, fr-fr, hu-hu, it-it, nl-nl, pl-pl, pt-pt, cs-cz, el-gr and zh-cn) and
// the maximum number of events (default: 10)
func (c *AuthClient) EventsCustomized(language string, max int) ([]EventData, error) {
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
		//fmt.Println("error:", err)
		return jsonResult, err
	}

	//fmt.Println(string(b))

	client := http.Client{}

	request, err := http.NewRequest("POST", c.getUrl("/api/v1/events/latest"), bytes.NewBuffer(b))
	if err != nil {
		//fmt.Println(err)
		return jsonResult, err
	}
	request.Header.Add("Content-Type", "application/json")

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errReq := client.Do(request)
	if errReq != nil {
		//fmt.Println(errReq)
		return jsonResult, err
	}
	//fmt.Println(response)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//panic(err)
		return jsonResult, err
	}
	//sb := string(body)
	//fmt.Println("raw body output:")
	//fmt.Println(sb)

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		//fmt.Println("Error in json unmarshalling")
		//fmt.Println(errJson)
		return jsonResult, errJson

	}
	//fmt.Println(jsonResult)

	return jsonResult, nil

}

func (c *AuthClient) Events() ([]EventData, error) {
	jsonResult := []EventData{}

	client := http.Client{}

	request, err := http.NewRequest("GET", c.getUrl("/api/v1/events/latest"), nil)
	if err != nil {
		return jsonResult, err
	}

	request.Header.Add("authorization", "Session "+c.SessionId)

	response, errMe := client.Do(request)
	if errMe != nil {
		return jsonResult, errMe
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonResult, err

	}
	//sb := string(body)
	//fmt.Println("raw body output:")
	//fmt.Println(sb)

	//fmt.Println(response.Body)
	//var resultMe map[string]interface{}
	//jsonResult := []EventData{}
	//jsonResult interface{}

	errJson := json.Unmarshal(body, &jsonResult)
	if errJson != nil {
		//fmt.Println("Error in json unmarshalling")
		//fmt.Println(errJson)
		return jsonResult, errJson

	}
	//fmt.Println(jsonResult)

	return jsonResult, nil

}
