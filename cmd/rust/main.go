package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/ohatakky/github-trending/pkg/trending"
	"github.com/ohatakky/github-trending/pkg/tweet"
)

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func health(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	if ip != "0.1.0.1" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(fmt.Sprintf("invalid IP: %s", ip))
		return
	}
	if r.Header.Get("X-Appengine-Cron") != "true" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid access: not appengine cron")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	cli := trending.New()
	twitter := tweet.New(os.Getenv("API_KEY_RUST"), os.Getenv("API_SECRET_RUST"), os.Getenv("ACCESS_TOKEN_RUST"), os.Getenv("ACCESS_TOKEN_SECRET_RUST"))

	job := func() {
		log.Println("-------------------- start --------------------")
		defer func() {
			log.Println("-------------------- end --------------------")
		}()

		items, err := cli.Daily("rust")
		if err != nil {
			log.Printf("fetch trending is failed: %s\n", err.Error())
			return
		}

		for _, item := range items {
			_, err := twitter.Tweet(fmt.Sprintf("%s %s", item.Text, item.Link))
			if err != nil {
				log.Printf("post tweet is failed: %s\n", err.Error())
				continue
			}

			time.Sleep(5 * time.Minute)
		}
	}

	scheduler.Every().Day().At("2:30").Run(job)

	// health check to avoid idle-timeout
	http.HandleFunc("/health", health)
	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
