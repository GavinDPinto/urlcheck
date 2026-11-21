package main

import (
	"bufio"
	"flag"
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
	var filePath string
	var numWorkers int
	var timeout time.Duration
	flag.StringVar(&filePath, "file", "urls.txt", "relative path to file with urls")
	flag.IntVar(&numWorkers, "workers", 5, "size of worker pool")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "timeout for GET requests")
	flag.Parse()

	if numWorkers < 1 {
		log.Fatal("-workers must be at least 1")
	}

	urls, err := readURLs(filePath)
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan string)
	results := make(chan Result)
	var wg sync.WaitGroup
	client := &http.Client{Timeout: timeout}

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
	down := 0
	for result := range results {
		if result.err != nil {
			fmt.Printf("%-5s %-4s %8v  %s %s\n", "DOWN", "-", result.duration.Round(time.Millisecond), result.url, result.err)
			down++
		} else if result.status >= 200 && result.status <= 399 {
			fmt.Printf("%-5s %-4d %8v  %s\n", "UP", result.status, result.duration.Round(time.Millisecond), result.url)
		} else {
			fmt.Printf("%-5s %-4d %8v  %s\n", "DOWN", result.status, result.duration.Round(time.Millisecond), result.url)
			down++
		}
	}
	duration := time.Since(start)

	fmt.Printf("Checking all %d URLs with %d worker(s) took a total time of %s.\n", len(urls), numWorkers, duration)
	fmt.Printf("Summary: %d of %d down.\n", down, len(urls))

	if down > 0 {
		os.Exit(1)
	}
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
