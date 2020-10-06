package function

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ohatakky/github-trending/pkg/trending"
	"github.com/ohatakky/github-trending/pkg/tweet"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cli := trending.New()
	items, err := cli.Read()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	twitter := tweet.New(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
	for _, item := range items {
		err := twitter.Tweet(item.Link)
		if err != nil {
			log.Println(err)
			continue
		}
		time.Sleep(2 * time.Minute)
	}
}
