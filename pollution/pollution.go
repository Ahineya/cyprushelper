package pollution

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"time"
	"strings"
	"errors"
	"strconv"
	"fmt"
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

// TODO: Refactor.
func GetPollution(city string) (*PollutionResult, error) {
	// TODO: this API updates every hour, but not in :00 minutes. There is a gap to :15 minutes now. Need to fix.
	from := time.Now().Add(-time.Hour).Format("02/01/2006+15:04")
	to := time.Now().Format("02/01/2006+15:04")

	city = normalizeCity(city)
	if city == "" {
		return nil, errors.New("There is no info for this city")
	}

	requestStr := `ajax_call=1&post_data[action]=get_plot_values_sd_dates`
	requestStr += `&post_data[station_code]=` + city
	requestStr += `&post_data[date_from]=`+ from
	requestStr += `&post_data[date_to]=`+ to
	requestStr += `&post_data[db_period]=1-hour`

	request, err := url.Parse(requestStr)
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
		return nil, errors.New(`There is no pollution data for this time, you can repeat the request later.
					Usually this happens at the beginning of an hour.`)
	}

	return &pollutionResult, nil
}

func FormatPollution(pollution *PollutionResult) string {
	result := ""

	for _, pollutionData := range pollution.Data {
		result += getPollutionLevel(pollutionData.PollutantCode, pollutionData.Value) + " "
		result += normalizePollutantCode(pollutionData.PollutantCode) + ": "
		result += "<b>" + pollutionData.Value + "</b> μg/m³\n"
	}

	return result
}

func normalizeCity(city string) string {
	if city == "limassol" || city == "lemesos" {
		return "LIMTRA"
	}

	if city == "pafos" || city == "paphos" {
		return "PAFTRA"
	}

	if city == "larnaka" || city == "larnaca" {
		return "LARTRA"
	}

	if city == "lefcosia" || city == "nicosia" {
		return "NICTRA"
	}

	return ""
}



func normalizePollutantCode(pollutantCode string) string {
	pollutants := make(map[string]string)
	pollutants["NO"] = "Nitrogen Oxide"
	pollutants["NO2"] = "Nitrogen Dioxide"
	pollutants["NOx"] = "Nitrogen Oxides"
	pollutants["SO2"] = "Sulfur Dioxide"
	pollutants["O3"] = "Ozone"
	pollutants["CO"] = "Carbon Monoxide"
	pollutants["PM10"] = "Particulate Matter 10 μm"
	pollutants["PM25"] = "Particulate Matter 2.5 μm"
	pollutants["BEN"] = "Benzene"

	return pollutants[pollutantCode]
}

func getPollutionLevel(pollutantCode string, value string) string {
	intValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	if pollutantCode == "PM10" {
		if intValue < 50 {
			return "[<b>Low</b>]"
		} else if intValue >= 50 && intValue < 100 {
			return "[<b>Moderate</b>]"
		} else if intValue >= 100 && intValue < 200 {
			return "[<b>High</b>]"
		} else if intValue >= 200 {
			return "[<b>Very high</b>]"
		}
	}
	return ""
}