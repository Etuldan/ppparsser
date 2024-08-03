package main

import (
	"time"

	"golang.org/x/net/html"
)

type Header struct {
	Title   string
	Date    time.Time
	Article []Article
	Id      string
}

func NewHeader(n *html.Node) *Header {
	header := &Header{}
	header.getAllElements(n)
	header.Date = time.Now()
	return header
}

type fn func(n *html.Node)

func (h *Header) SetTitle(n *html.Node, key string, val string) {
	for _, p := range n.Attr {
		if h.Title == "" && p.Key == key && p.Val == val {
			h.Title = n.FirstChild.Data
			return
		}
	}
}

func (h *Header) ExecuteOnTarget(n *html.Node, key string, val string, f fn) {
	for _, p := range n.Attr {
		if p.Key == key && p.Val == val {
			f(n)
			return
		}
	}
}

func (h *Header) AddArticle(n *html.Node) {
	article, err := NewArticle(n)
	if err != nil {
		return
	}
	h.Article = append(h.Article, *article)
}

func (h *Header) getAllElements(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		h.ExecuteOnTarget(n, "class", "sign-carousel--item sign-carousel--item--active", h.AddArticle)
		h.ExecuteOnTarget(n, "class", "sign-carousel--item ", h.AddArticle)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		h.getAllElements(c)
	}
}
