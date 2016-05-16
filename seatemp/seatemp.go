package seatemp

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
)

const (
	seatemp_url = "http://www.seatemperature.org/europe/cyprus/limassol.htm"
	seatemp_selector = "#sea-temperature"
)

func GetSeaTemp() (string, error){
	doc, err := goquery.NewDocument(seatemp_url)

	if err != nil {
		return "", err
	}

	seatemp, err := getSeatemp(doc)
	if err != nil {
		return "", err
	}

	return seatemp, nil
}

func FormatSeatemp(seatemp string) string {
	return seatemp
}

func getSeatemp(doc *goquery.Document) (string, error) {

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