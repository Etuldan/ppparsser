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
	a.ExecuteOnTarget(n, "span", "class", "date", a.SetDate)
	a.ExecuteOnTarget(n, "p", "class", "city", a.SetLocation)
	a.ExecuteOnTarget(n, "div", "class", "title", a.SetTitle)
	a.ExecuteOnTarget(n, "div", "class", "content", a.SetContent)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.getAllElements(c)
	}
}

func (a *Article) ExecuteOnTarget(n *html.Node, el string, key string, val string, f fnArticle) {
	if n.Type == html.ElementNode && n.Data == el {
		for _, p := range n.Attr {
			if p.Key == key && p.Val == val {
				f(n)
				return
			}
		}
	}
}

type fnArticle func(n *html.Node)

func (a *Article) SetContent(n *html.Node) {
	if len(a.Message) == 0 {
		var bufInnerHtml bytes.Buffer
		w := io.Writer(&bufInnerHtml)
		html.Render(w, n)
		a.Message = strings.Replace(bufInnerHtml.String(), "color: white;", "", -1)
		bufInnerHtml.Reset()
	}
}

func (a *Article) SetLocation(n *html.Node) {
	if len(a.Location) == 0 {
		a.Location = strings.Trim(n.FirstChild.Data, " \n")
	}
}
func (a *Article) SetTitle(n *html.Node) {
	if len(a.Title) == 0 {
		a.Title = strings.Trim(n.FirstChild.Data, " \n")
	}
}

func (a *Article) SetDate(n *html.Node) {
	if a.Date.IsZero() {
		regex := regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4})`)
		res := regex.FindAllStringSubmatch(n.FirstChild.Data, -1)
		if len(res) > 0 && len(res[0]) > 0 {
			date, _ := time.Parse("02/01/2006", res[0][0])
			a.Date = time.Date(date.UTC().Year(), date.UTC().Month(), date.UTC().Day(), 12, 0, 0, 0, date.UTC().Location())
		} else {
			now := time.Now().UTC()
			a.Date = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
		}
	}
}
