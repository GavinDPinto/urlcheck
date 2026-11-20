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
	file, err := os.Open("urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	test1 := check("https://go.dev")
	fmt.Println(test1.url, test1.status, test1.duration)

	test2 := check("https://does-not-exist.invalid")
	fmt.Println(test2.url, test2.status, test2.duration)
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
