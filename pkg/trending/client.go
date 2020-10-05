package trending

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const endpoint = "https://github.com/trending"

type Client struct{}

func New() *Client {
	return &Client{}
}

type Item struct {
	Link string
}

func (*Client) Read() ([]*Item, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	articles := doc.Find("article")
	items := make([]*Item, 0, articles.Length())
	articles.Each(func(i int, s *goquery.Selection) {
		h1 := s.Find("h1").First()
		a := h1.Find("a").First()
		link, exist := a.Attr("href")
		if !exist {
			return
		}
		items = append(items, &Item{
			Link: "https://github.com" + link,
		})
	})

	return items, nil
}
