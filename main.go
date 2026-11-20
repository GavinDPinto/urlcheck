package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

	start := time.Now()
	for _, url := range urls {
		result := check(url)
		if result.err != nil {
			fmt.Println(result.url, "ERROR:", result.err)
		} else {
			fmt.Println(result.url, result.status, result.duration)
		}
	}
	duration := time.Since(start)

	fmt.Println("Checking all URLs one-by-one took a total time of:", duration)
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

func check(url string) Result {
	client := &http.Client{Timeout: 5 * time.Second}

	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)
	if err != nil {
		return Result{url: url, duration: duration, err: err}
	}
	defer resp.Body.Close()

	return Result{url: url, duration: duration, status: resp.StatusCode}

}
