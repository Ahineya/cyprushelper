package pharmacies

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/kennygrant/sanitize"
	"net/http"
	"io/ioutil"
	"bytes"
	"fmt"
	"strconv"
	"github.com/Ahineya/cyprushelper/utils"
)

const (
	pharmacies_url = "http://www.cypruspharmacy.com/"
	pharmacies_cities_selector = ".featuredHeader"
	pharmacies_names_selector = ".middle_col > table:nth-child(%d) .pharmacyheader:nth-child(5n+2)"
	pharmacies_phones_selector = ".middle_col table:nth-child(%d) .pharmacyheader:nth-child(5n+3)"
	pharmacies_home_phones_selector = ".middle_col table:nth-child(%d) .pharmacyheader:nth-child(5n+4)"
	pharmacies_addresses_selector = ".middle_col > table:nth-child(%d) .newstitle~td"
)

type Pharmacy struct {
	Name      string
	Address   string
	Phone     string
	HomePhone string
}

type PharmaciesInCity struct {
	City       string
	Pharmacies []Pharmacy
}

type Cities []string

type AllPharmacies []PharmaciesInCity

func GetAllPharmacies() (AllPharmacies, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", pharmacies_url, nil)
	if err != nil {
		return AllPharmacies{}, err
	}

	req.Close = true

	req.Header.Add("Accept-Encoding", "identity")

	resp, err := client.Do(req)

	if err != nil {
		return AllPharmacies{}, err
	}

	defer resp.Body.Close()

	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AllPharmacies{}, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))

	if err != nil {
		return AllPharmacies{}, err
	}

	cities, err := getCities(doc)
	if err != nil {
		return AllPharmacies{}, err
	}

	var allPharmacies AllPharmacies

	for cityId, city := range cities {
		pharmacies, err := getPharmaciesByCityId(doc, cityId)
		if err != nil {
			return AllPharmacies{}, err
		}

		allPharmacies = append(allPharmacies, PharmaciesInCity{city, pharmacies})
	}

	return allPharmacies, nil
}

func FormatPharmacies(pharmaciesInCity PharmaciesInCity) string {
	result := "<b>" + utils.UpcaseInitial(pharmaciesInCity.City) + "</b>\n\n"

	for idx, pharmacy := range pharmaciesInCity.Pharmacies {
		result += "<b>" + strconv.Itoa(idx + 1) + ". "
		result += "" + pharmacy.Name + "</b>\n"
		result += "Address: " + pharmacy.Address + "\n"
		result += "Phone: " + pharmacy.Phone + "\n"
		result += "Home Phone: " + pharmacy.HomePhone + "\n\n"
	}

	return result
}

func getCities(doc *goquery.Document) (Cities, error) {

	elems := doc.Find(pharmacies_cities_selector)
	elemsLength := elems.Length()

	var cities Cities

	if (elemsLength == 0) {
		return cities, errors.New("Can't find pharmacies in html")
	}

	for i := 0; i < elemsLength; i++ {
		pharmacyHtml, err := elems.Eq(i).Html()
		if err != nil {
			return cities, errors.New("Can't get pharmacy html")
		}
		pharmacy, err := sanitize.HTMLAllowing(pharmacyHtml, []string{"b", "strong", "i", "em", "a", "code", "br"})
		if err != nil {
			return cities, errors.New("Can't parse pharmacy")
		}
		pharmacy = utils.Replace(pharmacy, "<br/>", " ")

		cities = append(cities, pharmacy)
	}

	return cities, nil
}

func getPharmaciesByCityId(doc *goquery.Document, cityId int) ([]Pharmacy, error) {
	pharmacyNames, err := getPharmacyNames(doc, cityId)
	if err != nil {
		return []Pharmacy{}, err
	}

	pharmacyAddresses, err := getPharmacyAddresses(doc, cityId)
	if err != nil {
		return []Pharmacy{}, err
	}

	pharmacyPhones, err := getPharmacyPhones(doc, cityId)
	if err != nil {
		return []Pharmacy{}, err
	}

	pharmacyHomePhones, err := getPharmacyHomePhones(doc, cityId)
	if err != nil {
		return []Pharmacy{}, err
	}

	var pharmacies []Pharmacy

	for idx, _ := range pharmacyNames {
		pharmacy := Pharmacy{pharmacyNames[idx], pharmacyAddresses[idx], pharmacyPhones[idx], pharmacyHomePhones[idx]}
		pharmacies = append(pharmacies, pharmacy)
	}

	return pharmacies, nil
}

func getPharmacyNames(doc *goquery.Document, cityId int) ([]string, error) {
	return getDataForCity(doc, cityId, pharmacies_names_selector)
}

func getPharmacyPhones(doc *goquery.Document, cityId int) ([]string, error) {
	return getDataForCity(doc, cityId, pharmacies_phones_selector)
}

func getPharmacyHomePhones(doc *goquery.Document, cityId int) ([]string, error) {
	return getDataForCity(doc, cityId, pharmacies_home_phones_selector)
}

func getPharmacyAddresses(doc *goquery.Document, cityId int) ([]string, error) {
	return getDataForCity(doc, cityId, pharmacies_addresses_selector)
}

func getDataForCity(doc *goquery.Document, cityId int, selector string) ([]string, error) {

	elems := doc.Find(fmt.Sprintf(selector, cityId * 2 + 3))
	elemsLength := elems.Length()

	var data []string

	if (elemsLength == 0) {
		return data, errors.New("Can't find pharmacies in html")
	}

	for i := 0; i < elemsLength; i++ {
		pharmacyHtml, err := elems.Eq(i).Html()
		if err != nil {
			return data, errors.New("Can't get pharmacy html")
		}
		pharmacy, err := sanitize.HTMLAllowing(pharmacyHtml, []string{"b", "strong", "i", "em", "a", "code", "br"})
		if err != nil {
			return data, errors.New("Can't parse pharmacy")
		}
		pharmacy = utils.Replace(pharmacy, "<br/>", " ")

		data = append(data, pharmacy)
	}

	return data, nil
}

func NormalizeCity(city string) string {
	if city == "limassol" || city == "lemesos" {
		return "limassol"
	}

	if city == "pafos" || city == "paphos" {
		return "paphos"
	}

	if city == "larnaka" || city == "larnaca" {
		return "larnaca"
	}

	if city == "lefcosia" || city == "nicosia" {
		return "nicosia"
	}

	if city == "famagusta" {
		return "famagusta"
	}

	return city
}