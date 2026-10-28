
package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

func main() {
	urls := []string{
		"https://go.dev",
		"https://github.com",
		"https://example.com",
		"https://golang.org",
		"https://httpbin.org/status/404",
	}

	const workerCount = 3

	jobs := make(chan string)
	results := make(chan Result)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup

	// Start a fixed number of workers.
	wg.Add(workerCount)

	for id := 1; id <= workerCount; id++ {
		go worker(id, jobs, results, client, &wg)
	}

	// Producer: send URLs and close the jobs channel.
	go func() {
		defer close(jobs)

		for _, url := range urls {
			jobs <- url
		}
	}()

	// Coordinator: close results after every worker exits.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Consumer: collect results until the channel is closed.
	for result := range results {
		if result.Err != nil {
			fmt.Printf(
				"ERROR      %-28s %v\n",
				result.URL,
				result.Err,
			)
			continue
		}

		label := "HEALTHY"
		if result.StatusCode >= 400 {
			label = "UNHEALTHY"
		}

		fmt.Printf(
			"%-10s HTTP %d  %-28s %v\n",
			label,
			result.StatusCode,
			result.URL,
			result.Duration.Round(time.Millisecond),
		)
	}

	fmt.Println("\nAll website checks completed.")
}

func worker(
	id int,
	jobs <-chan string,
	results chan<- Result,
	client *http.Client,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for url := range jobs {
		fmt.Printf("Worker %d checking %s\n", id, url)

		start := time.Now()
		response, err := client.Get(url)

		result := Result{
			URL:      url,
			Duration: time.Since(start),
			Err:      err,
		}

		if err == nil {
			result.StatusCode = response.StatusCode
			response.Body.Close()
		}

		results <- result
	}
}
