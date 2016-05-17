package pollution

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"time"
	"strings"
	"errors"
)

const pollution_url = "http://www.airquality.dli.mlsi.gov.cy/site/ajax_exec"

type PollutionData struct {
	DateTime          string `json:"date_time"`
	IsPollutant       bool   `json:"is_pollutant"`
	PollutantCode     string `json:"pollutant_code"`
	PollutantHTMLCode string `json:"pollutant_html_code"`
	Uom               string `json:"uom"`
	Value             string `json:"value"`
}

type PollutionResult struct {
	Data []PollutionData `json:"data"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

// TODO: Refactor. Add cities.
func GetPollution(city string) (*PollutionData, error) {
	from := time.Now().Add(-time.Hour).Format("02/01/2006+15:04")
	to := time.Now().Format("02/01/2006+15:04")

	request, err := url.Parse(`ajax_call=1&post_data[action]=get_plot_values_sd_dates&post_data[station_code]=LIMTRA&post_data[date_from]=`+ from +`&post_data[date_to]=`+ to +`&post_data[db_period]=1-hour`)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", pollution_url, strings.NewReader(request.String()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pollutionResult PollutionResult;
	err = json.Unmarshal(body, &pollutionResult)
	if err != nil {
		return nil, err
	}

	if len(pollutionResult.Data) == 0 {
		return nil, errors.New("There is no pollution data for this time and city")
	}

	for _, data := range pollutionResult.Data {
		if data.PollutantCode == "PM10" {
			return &data, nil
		}
	}

	return nil, errors.New("There is no pollution data for this time and city")
}

func FormatPollution(pollution *PollutionData) string {
	result := ""
	result += "<b>" + pollution.PollutantCode + "</b>: " + pollution.Value
	return result
}
