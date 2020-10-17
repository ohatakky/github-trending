package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/ohatakky/github-trending/pkg/trending"
	"github.com/ohatakky/github-trending/pkg/tweet"
)

func main() {
	cli := trending.New()
	twitter := tweet.New(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))

	job := func() {
		log.Println("-------------------- start --------------------")
		defer func() {
			log.Println("-------------------- end --------------------")
		}()

		items, err := cli.Read()
		if err != nil {
			log.Printf("fetch trending is failed: %s\n", err.Error())
			return
		}

		for _, item := range items {
			_, err := twitter.Tweet(item.Link)
			if err != nil {
				log.Printf("post tweet is failed: %s\n", err.Error())
				continue
			}

			time.Sleep(5 * time.Minute)
		}
	}

	scheduler.Every().Day().At("1:10").Run(job)
	runtime.Goexit()
}
