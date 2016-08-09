package main

import (
	"bufio"
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	// Open file
	var filename = flag.String("file", "urls.txt", "File to check")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer file.Close()

	// Read File

	lines, err := readFile(file)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(lines)

}

func fetch(link string, ch chan<- string) (*http.Response, error) {
	resp, err := http.Head(link)
	if err != nil {
		log.Fatal(err)
	}
	return resp, err
}

// readFile reads in a file and returns a slice of strings with the spaces removed
func readFile(file *os.File) ([]string, error) {
	name := file.Name()
	log.Printf("Reading file: %s \n", name)

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()

		if err := scanner.Err(); err != nil {
			return lines, err
		}

		line = strings.TrimSpace(line)
		lines = append(lines, line)
	}

	return lines, nil
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}
