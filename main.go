package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"golang.org/x/net/html"
)

func createFeed() string {
	city := 890589053 // FIXME
	requestURL := fmt.Sprintf("https://app.panneaupocket.com/embeded/%d?mode=widgetTv", city)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	header := NewHeader(doc)

	feed := &feeds.Feed{
		Title:       header.Title,
		Description: "Actualités de " + fmt.Sprint(city),
		Link:        &feeds.Link{Href: "https://info-communes.fr"},
		Created:     time.Now().UTC(),
		Id:          fmt.Sprint(city),
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
		log.Fatal(err)
	}
	return rss
}

func renderFeed(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(createFeed()))
}

func main() {
	http.HandleFunc("/", renderFeed)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
