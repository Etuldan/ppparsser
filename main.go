package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"golang.org/x/net/html"
)

func createFeed(city *City) (string, error) {
	requestURL := fmt.Sprintf("https://app.panneaupocket.com/embeded/%d?mode=widgetTv", city.Id)
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", fmt.Errorf("error while accessing embed api")
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading embed api")
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
		return "", fmt.Errorf("error while creating rss")
	}
	return rss, nil
}

func renderFeed(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}

	city, err := NewCity(r.URL.Path[1:])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}
	err = city.GetIdFromPP()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}
	rss, err := createFeed(city)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "https://info-communes.fr", http.StatusSeeOther)
		return
	}
	w.Write([]byte(rss))
}

func main() {
	http.HandleFunc("/", renderFeed)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
