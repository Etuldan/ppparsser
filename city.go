package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type City struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Postcode string `json:"postcode"`
}

func NewCity(insee string) (*City, error) {
	city := &City{}
	name, err := city.queryLaposte("nom_de_la_commune", insee)
	if err != nil {
		return nil, err
	}
	city.Name = name
	postcode, err := city.queryLaposte("code_postal", insee)
	if err != nil {
		return nil, err
	}
	city.Postcode = postcode
	return city, nil
}

func (c *City) GetIdFromPP() error {
	resp, err := http.Get("https://app.panneaupocket.com/public-api/city")
	if err != nil {
		return fmt.Errorf("error while accessing city api")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading city api")
	}
	var cities []City
	json.Unmarshal([]byte(body), &cities)
	defer resp.Body.Close()

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	for _, city := range cities {
		if city.Postcode == c.Postcode {
			normalized, _, err := transform.String(t, strings.Replace(city.Name, "-", " ", -1))
			if err != nil {
				return err
			}
			if strings.EqualFold(normalized, c.Name) {
				c.Id = city.Id
				return nil
			}
		}
	}
	return fmt.Errorf("city not found")
}

func (c *City) queryLaposte(field string, value string) (string, error) {
	requestUrl := fmt.Sprintf("https://datanova.laposte.fr/data-fair/api/v1/datasets/laposte-hexasmal/values/%s?size=1&q_fields=code_commune_insee&qs=%s", field, value)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return "", fmt.Errorf("error while accessing laposte api")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading laposte api")
	}
	var cities []string
	json.Unmarshal([]byte(body), &cities)
	if len(cities) == 0 {
		return "", fmt.Errorf("invalid insee")
	}
	return cities[0], nil
}
