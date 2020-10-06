package main

import (
	_ "expvar"
	"log"
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
	cli := trending.New()
	items, err := cli.Read()
	if err != nil {
		log.Fatal(err)
	}

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
}
