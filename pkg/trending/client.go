package trending

import (
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"
)

const endpoint = "https://github.com/trending"

type Client struct{}

func New() *Client {
	return &Client{}
}

type Item struct {
	Link string
	Text string
}

func (*Client) Daily(lang string) ([]*Item, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if lang != "" {
		u.Path = path.Join(u.Path, lang)
	}
	q := u.Query()
	q.Set("since", "daily")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	articles := doc.Find("article")
	items := make([]*Item, 0, articles.Length())
	articles.Each(func(i int, s *goquery.Selection) {
		h1 := s.Find("h1").First()
		a := h1.Find("a").First()
		link, exist := a.Attr("href")
		if !exist {
			log.Println("the html structure has changed")
			return
		}
		text := s.Find("p").First().Text()
		items = append(items, &Item{
			Link: "https://github.com" + link,
			Text: text,
		})
	})

	return items, nil
}
