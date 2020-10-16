package main

import (
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/ohatakky/github-trending/pkg/trending"
	"github.com/ohatakky/github-trending/pkg/tweet"
)

func main() {
	job := func() {
		cli := trending.New()
		items, err := cli.Read()
		if err != nil {
			log.Fatalf("fetch trending is failed: %s\n", err.Error())
		}

		queue := make(chan trending.Item)
		twitter := tweet.New(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range queue {
				if twitter.Tweet(item.Link) != nil {
					log.Printf("post tweet is failed: %s\n", err.Error())
				}
				time.Sleep(4 * time.Minute)
			}
		}()

		for _, item := range items {
			queue <- *item
		}

		close(queue)
		wg.Wait()
	}

	scheduler.Every().Day().At("08:00").Run(job)

	runtime.Goexit()
}
