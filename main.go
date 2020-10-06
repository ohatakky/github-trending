package main

import (
	_ "expvar"
	"fmt"
	"sync"
	"time"
)

const worker = 1

const duration = 3

type job struct {
	name string
}

func doWork(j job) {
	time.Sleep(duration * time.Second)
	fmt.Println(j.name)
}

func main() {
	jobs := make(chan job)

	// start worker
	wg := &sync.WaitGroup{}
	wg.Add(worker)

	go func() {
		defer wg.Done()
		for j := range jobs {
			doWork(j)
		}
	}()

	// add jobs
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("job-%d", i)
		jobs <- job{name}
	}

	close(jobs)
	wg.Wait()
}
