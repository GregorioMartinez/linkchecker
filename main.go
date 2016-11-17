package main

import (
	"bufio"
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	// Open file
	var filename = flag.String("file", "", "File to check")
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

	for _, link := range lines {
		resp, err := fetch(link)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// @TODO test this
		robots := resp.Header.Get("X-Robots-Tag")

		fmt.Println(resp.Status, link, robots)

	}
}

func fetch(link string) (*http.Response, error) {
	resp, err := http.Head(link)
	if err != nil {
		// log.Fatal(err)
	}
	return resp, err
}

// readFile reads in a file and returns a slice of strings with the spaces removed
// @TODO make sure valid url
// @TODO make relative to absolute
// Reads in txt files right now. Expand to txt, csv, xml
// @TODO don't just trust file extension
func readFile(file *os.File) ([]string, error) {
	name := file.Name()

	log.Printf("Reading file: %s \n", name)

	extension := path.Ext(name)

	switch extension {
	case ".txt":
		fmt.Println("txt file loaded")
	case ".xml":
		fmt.Println("XML document loaded")
	default:
		fmt.Println("unsupported file extension")
	}

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
