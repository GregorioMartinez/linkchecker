package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var requestClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func main() {
	// Open file
	var filename = flag.String("file", "", "File to pull urls from. Accepts .txt, .csv, .xml")
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

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			fmt.Println(resp.Request.URL.String())
		}

		// @TODO test this
		robots := resp.Header.Get("X-Robots-Tag")

		fmt.Println(resp.Status, link, robots)

	}
}

func fetch(link string) (*http.Response, error) {
	resp, err := requestClient.Head(link)
	if err != nil {
		log.Fatal(err)
	}
	return resp, err
}

type ListLinks interface {
	getLinks(file *os.File) ([]string, error)
}

type TxtReader struct {
}

type XmlDoc struct {
	Locs []string `xml:"url>loc"`
}

type XmlReader struct{}

//@TODO Should i get a byte array rather than pas around a file?
func (r XmlReader) getLinks(file *os.File) ([]string, error) {

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return []string{}, err
	}

	var doc XmlDoc
	err = xml.Unmarshal(b, &doc)
	if err != nil {
		return []string{}, err
	}
	return doc.Locs, err
}

func (r TxtReader) getLinks(file *os.File) ([]string, error) {
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
		tr := TxtReader{}
		return tr.getLinks(file)
	case ".xml":
		xr := XmlReader{}
		return xr.getLinks(file)
	case ".csv":

	default:
		fmt.Println("unsupported file extension")
	}

	return []string{}, nil
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}
