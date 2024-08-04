package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	Insee    string
}

func NewCity(insee string) *City {
	city := &City{Insee: insee}

	return city
}

func (c *City) Populate(local bool) error {
	if local {
		err := c.readCSV()
		if err != nil {
			return err
		}
	} else {
		name, err := c.queryLaposte("nom_de_la_commune")
		if err != nil {
			return err
		}
		c.Name = name
		postcode, err := c.queryLaposte("code_postal")
		if err != nil {
			return err
		}
		c.Postcode = postcode
	}
	return nil
}

func (c *City) readCSV() error {
	file, err := os.Open("./data/base-officielle-codes-postaux.csv")
	if err != nil {
		return fmt.Errorf("error while reading the file : %s", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading records : %s", err)
	}
	for _, line := range records {
		if line[0] == c.Insee {
			c.Name = line[1]
			c.Postcode = line[2]
			return nil
		}
	}
	return fmt.Errorf("no city found %s", c.Insee)
}

func (c *City) GetIdFromPP(local bool) error {
	var data []byte
	if local {
		file, err := os.ReadFile("./data/pp-city.json")
		if err != nil {
			return fmt.Errorf("error while reading the file : %s", err)
		}
		data = file
	} else {
		resp, err := http.Get("https://app.panneaupocket.com/public-api/city")
		if err != nil {
			return fmt.Errorf("error while accessing city api : %s", err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error while reading city api : %s", err)
		}
		data = []byte(body)
		defer resp.Body.Close()
	}

	var cities []City
	json.Unmarshal(data, &cities)

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
	return fmt.Errorf("city not found %s %s", c.Postcode, c.Name)
}

func (c *City) queryLaposte(field string) (string, error) {
	requestUrl := fmt.Sprintf("https://datanova.laposte.fr/data-fair/api/v1/datasets/laposte-hexasmal/values/%s?size=1&q_fields=code_commune_insee&qs=%s", field, c.Insee)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return "", fmt.Errorf("error while accessing laposte api : %s", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading laposte api : %s", err)
	}
	var cities []string
	json.Unmarshal([]byte(body), &cities)
	if len(cities) == 0 {
		return "", fmt.Errorf("no data found on laposte api for %s %s", field, c.Insee)
	}
	return cities[0], nil
}
