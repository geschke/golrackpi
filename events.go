package golrackpi

import (
	"bytes"
	"encoding/json"

	"io/ioutil"
	"net/http"

	//"time"
	"github.com/geschke/golrackpi/internal/timefix"
)

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
