package pharmacies

import (
	"errors"
	"regexp"
	"github.com/PuerkitoBio/goquery"
	"github.com/kennygrant/sanitize"
)

const (
	pharmacies_url = "https://www.cyta.com.cy/id/m144/en"
	pharmacies_selector = "#ctl00_cphContent_grvPharmacies td"
)

/*
	 We are getting Cyta page with night pharmacies here. By default it will get pharmacies for Nicosia
	 To get another cities pharmacies, we need to bypass the Cyta security. I have tried to take the curl
	  request from the chrome network page, and it works from command line. But it doesn't work from Postman.
	  I think we should take a look to cookies that Cyta site sets. (Refer to the bypassASP.GetASPViewStateVars)
 */
func GetPharmacies() ([]string, error) {
	doc, err := goquery.NewDocument(pharmacies_url)
	if err != nil {
		return []string{}, err
	}

	pharmacies, err := getPharmacies(doc)
	if err != nil {
		return []string{}, err
	}

	return pharmacies, nil
}

func getPharmacies(doc *goquery.Document) ([]string, error) {
	elems := doc.Find(pharmacies_selector)
	elemsLength := elems.Length()

	var pharmacies []string

	if (elemsLength == 0) {
		return pharmacies, errors.New("Can't find pharmacies in html")
	}

	for i := 0; i < elemsLength; i++ {
		pharmacyHtml, err := elems.Eq(i).Html()
		if err != nil {
			return pharmacies, errors.New("Can't get pharmacy html")
		}
		pharmacy, err := sanitize.HTMLAllowing(pharmacyHtml, []string{"b", "strong", "i", "em", "a", "code", "br"})
		if err != nil {
			return pharmacies, errors.New("Can't parse pharmacy")
		}
		pharmacy = replace(pharmacy, "<br/>", " ")

		pharmacies = append(pharmacies, pharmacy)
	}

	return pharmacies, nil
}

func replace(str string, tag string, replacer string) string {
	r := regexp.MustCompile(tag)
	return r.ReplaceAllString(str, replacer)
}
