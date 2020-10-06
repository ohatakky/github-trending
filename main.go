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
	duration = 3
)

var twitter = tweet.New(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))

func work(job trending.Item) {
	time.Sleep(duration * time.Second)
	err := twitter.Tweet(job.Link)
	if err != nil {
		log.Println(err)
	}
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
			work(j)
		}
	}()

	for _, item := range items {
		jobs <- *item
	}

	close(jobs)
	wg.Wait()
}
