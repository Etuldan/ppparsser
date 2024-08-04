package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/feeds"
	"golang.org/x/net/html"
)

var regexInsee = regexp.MustCompile(`^\/(\d{5})$`)
var regexPP = regexp.MustCompile(`^\/panneaupocket\/(\d*)$`)

func createFeed(city *City) (string, error) {
	requestURL := fmt.Sprintf("https://app.panneaupocket.com/embeded/%d?mode=widgetTv", city.Id)
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", fmt.Errorf("error while accessing embed api : %s", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading embed api : %s", err)
	}

	header := NewHeader(doc)

	feed := &feeds.Feed{
		Title:       header.Title,
		Description: "Actualités de " + city.Name,
		Link:        &feeds.Link{Href: "https://info-communes.fr"},
		Created:     time.Now().UTC(),
		Id:          fmt.Sprint(city.Id),
		Items:       []*feeds.Item{},
	}

	for _, s := range header.Article {
		item := &feeds.Item{
			Title:       s.Title,
			Description: s.Message,
			Created:     s.Date,
			Id:          feed.Id + "/" + s.Id,
			Author:      &feeds.Author{Name: s.Location},
		}
		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", fmt.Errorf("error while creating rss : %s", err)
	}
	return rss, nil
}

func renderGeneral(w http.ResponseWriter, r *http.Request) {
	insee := regexInsee.FindStringSubmatch(r.URL.Path)
	if insee == nil || len(insee) < 2 {
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}

	city := NewCity(insee[1])
	err := city.Populate(true)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}
	err = city.GetIdFromPP(true)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}
	renderFeed(w, r, city)
}

func renderPanneaupocket(w http.ResponseWriter, r *http.Request) {
	pp := regexPP.FindStringSubmatch(r.URL.Path)
	if pp == nil || len(pp) < 2 {
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}

	idPP, err := strconv.Atoi(pp[1])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}
	city := &City{Id: idPP}
	renderFeed(w, r, city)
}

func renderFeed(w http.ResponseWriter, r *http.Request, city *City) {
	rss, err := createFeed(city)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}
	w.Write([]byte(rss))
}

func main() {
	http.HandleFunc("/panneaupocket/", renderPanneaupocket)
	http.HandleFunc("/", renderGeneral)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
