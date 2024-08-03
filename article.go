package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Article struct {
	Title    string
	Date     time.Time
	Message  string
	Id       string
	Location string
}

func NewArticle(n *html.Node) (*Article, error) {
	article := &Article{}
	for _, p := range n.Attr {
		if p.Key == "data-id" {
			article.Id = p.Val
		}
	}
	article.getAllElements(n)

	if article.Title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	return article, nil
}

func (a *Article) getAllElements(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "span" {
		a.SetDate(n)
	} else if n.Type == html.ElementNode && n.Data == "div" {
		a.SetTitle(n)
		a.SetContent(n)
	} else if n.Type == html.ElementNode && n.Data == "p" {
		a.SetLocation(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.getAllElements(c)
	}
}

type fnArticle func(n *html.Node)

func (a *Article) getInfo(n *html.Node, key string, val string, f fnArticle) {
	for _, p := range n.Attr {
		if p.Key == key && p.Val == val {
			f(n)
			return
		}
	}
}

func (a *Article) SetContent(n *html.Node) {
	if len(a.Message) == 0 {
		a.getInfo(n, "class", "content", func(n *html.Node) {
			var bufInnerHtml bytes.Buffer
			w := io.Writer(&bufInnerHtml)
			html.Render(w, n)
			a.Message = strings.Replace(bufInnerHtml.String(), "color: white;", "", -1)
			bufInnerHtml.Reset()
		})
	}
}

func (a *Article) SetLocation(n *html.Node) {
	if len(a.Location) == 0 {
		a.getInfo(n, "class", "city", func(n *html.Node) {
			a.Location = strings.Trim(n.FirstChild.Data, " \n")
		})
	}
}
func (a *Article) SetTitle(n *html.Node) {
	if len(a.Title) == 0 {
		a.getInfo(n, "class", "title", func(n *html.Node) {
			a.Title = strings.Trim(n.FirstChild.Data, " \n")
		})
	}
}

func (a *Article) SetDate(n *html.Node) {
	if a.Date.IsZero() {
		a.getInfo(n, "class", "date", func(n *html.Node) {
			regex := regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4})`)
			res := regex.FindAllStringSubmatch(n.FirstChild.Data, -1)
			if len(res) > 0 && len(res[0]) > 0 {
				date, _ := time.Parse("02/01/2006", res[0][0])
				a.Date = time.Date(date.UTC().Year(), date.UTC().Month(), date.UTC().Day(), 12, 0, 0, 0, date.UTC().Location())
			} else { // Aujourd'hui
				now := time.Now().UTC()
				a.Date = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
			}
		})
	}
}
