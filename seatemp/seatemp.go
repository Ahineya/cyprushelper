package seatemp

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

const (
	seatemp_url = "http://www.seatemperature.org/europe/cyprus/%s.htm"
	seatemp_selector = "#sea-temperature"
)

func GetSeaTemp(city string) (string, error){
	city = normalizeCity(city)
	doc, err := goquery.NewDocument(fmt.Sprintf(seatemp_url, city))

	if err != nil {
		return "", err
	}

	seatemp, err := getSeaTemp(doc)
	if err != nil {
		return "", err
	}

	return seatemp, nil
}

func FormatSeatemp(seatemp string) string {
	return seatemp
}

func getSeaTemp(doc *goquery.Document) (string, error) {

	elems := doc.Find(seatemp_selector)
	elemsLength := elems.Length()

	if (elemsLength == 0) {
		return "", errors.New("Can't find sea temperature in html")
	}

	seatemp, err := elems.Eq(0).Html()
	if err != nil {
		return "", err
	}

	return seatemp, nil
}

func normalizeCity(city string) string {
	if city == "limassol" || city == "lemesos" {
		return "limassol"
	}

	if city == "pafos" || city == "paphos" {
		return "kissonerga"
	}

	if city == "larnaka" || city == "larnaca" {
		return "perivolia"
	}

	if city == "lefcosia" || city == "nicosia" {
		return "kato-pyrgos"
	}

	if city == "famagusta" {
		return "famagusta"
	}

	if city == "protaras" {
		return "protaras"
	}

	if city == "ayia-napa" || city == "ayia" {
		return "ayia-napa"
	}

	return ""
}
