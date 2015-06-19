package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
)

func runWorkers(numOfWorders int, queueCapacity int, wg *sync.WaitGroup) chan<- string {

	createWorker := func(workerId int, recvCh <-chan string) func() {
		return func() {
			for url := range recvCh {
				res, err := http.Get(url)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[%d]%s: %v\n", workerId, url, err)
					continue
				}
				defer res.Body.Close()

				b, err := ioutil.ReadAll(res.Body)
				if err == nil {
					fmt.Printf("[%d]%s: %s, bytes=%d\n", workerId, url, res.Status, len(b))
				}
			}
		}
	}

	ch := make(chan string, queueCapacity)
	for i := 0; i < numOfWorders; i++ {
		wg.Add(1)
		worker := createWorker(i, ch)
		go func() {
			defer wg.Done()
			worker()
		}()
	}

	return ch
}

func main() {

	cpus := runtime.NumCPU()
	fmt.Fprintf(os.Stderr, "cpus=%d\n", cpus)
	runtime.GOMAXPROCS(cpus)

	urls := []string{
		"http://www.yahoo.co.jp/notfound",
		"http://www.yahoo.co.ijp/",
		"http://www.yahoo.co.jp/",
		"http://www.google.co.jp/",
		"http://www.goo.ne.jp/",
	}

	var wg sync.WaitGroup
	ch := runWorkers(3, 10, &wg)
	for _, url := range urls {
		ch <- url
	}
	close(ch)

	fmt.Fprintf(os.Stderr, "Waiting for all workers done.\n")
	wg.Wait()
}
