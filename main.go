package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	url      string
	status   int
	duration time.Duration
	err      error
}

func main() {
	urls, err := readURLs("urls.txt")
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan string)
	results := make(chan Result)
	var wg sync.WaitGroup
	client := &http.Client{Timeout: 5 * time.Second}

	const numWorkers = 5

	start := time.Now()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(client, jobs, results, &wg)
	}
	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()
	go func() {
		wg.Wait()
		close(results)
	}()
	for result := range results {
		if result.err != nil {
			fmt.Println(result.url, "ERROR:", result.err)
		} else {
			fmt.Println(result.url, result.status, result.duration)
		}
	}
	duration := time.Since(start)

	fmt.Printf("Checking all %d URLs with %d workers took a total time of %s.", len(urls), numWorkers, duration)
}

func worker(client *http.Client, jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		results <- check(url, client)
	}
}

func readURLs(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("reading URL list: %w", err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" && !strings.HasPrefix(url, "#") {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading URL list: %w", err)
	}

	return urls, nil
}

func check(url string, client *http.Client) Result {
	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)
	if err != nil {
		return Result{url: url, duration: duration, err: err}
	}
	defer resp.Body.Close()

	return Result{url: url, duration: duration, status: resp.StatusCode}

}
