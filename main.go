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

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func handler(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	if ip != "0.1.0.1" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(fmt.Sprintf("invalid IP: %s", ip))
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
		log.Printf("fetch trending is failed: %s\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

	jobs := make(chan trending.Item)
	twitter := tweet.New(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := range jobs {
			time.Sleep(10 * time.Minute)
			if twitter.Tweet(j.Link) != nil {
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

func main() {
	http.HandleFunc("/", handler)
	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
