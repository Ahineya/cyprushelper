package bypassASP

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"errors"
)

type ASPVars struct {
	VIEWSTATE string
	EVENTVALIDATION string
}

func GetASPViewStateVars(url string) (ASPVars, error) {
	//Going through http.Get because we need to read Set-Cookie header ( .CMP_SESSIONID cookie)
	resp, err := http.Get(url)
	if err != nil {
		return ASPVars{"",""}, err
	}
	//fmt.Println(resp.Header)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ASPVars{"",""}, err
	}

	viewstate, exists := doc.Find("#__VIEWSTATE").Attr("value")
	if !exists {
		return ASPVars{"",""}, errors.New("VIEWSTATE var not exists")
	}

	eventvalidation, exists := doc.Find("#__EVENTVALIDATION").Attr("value")
	if !exists {
		return ASPVars{"",""}, errors.New("EVENTVALIDATION var not exists")
	}

	// Encoding viewstate and eventvalidation for the request
	/*data := url.Values{}
	data.Set("__VIEWSTATE", viewstate)
	data.Add("__EVENTVALIDATION", eventvalidation)

	fmt.Println(data.Encode())*/

	return ASPVars{viewstate, eventvalidation}, nil
}