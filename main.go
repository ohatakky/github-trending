package main

import (
	_ "expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ohatakky/github-trending/pkg/trending"
	"github.com/ohatakky/github-trending/pkg/tweet"
)

const (
	worker   = 1
	duration = 10
)

var twitter = tweet.New(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))

func work(job trending.Item) error {
	time.Sleep(duration * time.Minute)
	return twitter.Tweet(job.Link)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "github-trending-dot-akki-256705.appspot.com" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(fmt.Sprintf("invalid host: %s", r.Host))
			return
		}
		if r.Header.Get("X-Appengine-Cron") != "true" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(fmt.Sprintf("invalid access: not appengine cron"))
			return
		}

		cli := trending.New()
		items, err := cli.Read()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)

		wg := &sync.WaitGroup{}
		wg.Add(worker)
		jobs := make(chan trending.Item)
		go func() {
			defer wg.Done()
			for j := range jobs {
				if work(j) != nil {
					log.Println(err)
				}
			}
		}()

		for _, item := range items {
			jobs <- *item
		}

		close(jobs)
		wg.Wait()
	})

	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
