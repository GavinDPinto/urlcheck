package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Milestone 1: read urls.txt and print each URL.
	file, err := os.Open("urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			fmt.Println(line)
		}
	}
	if err:= scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
