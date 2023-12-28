package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ScrapConfig struct {
	ReceiverCSV string `json:"receiverCSV"`
	ScrapURL    string `json:"scrapURL"`
	HeaderCSV   string `json:"headersCSV"`
	// SearchWord  string `json:"searchWord"`
}

type Scraper interface {
	scrapCFG() string
}

func (s *ScrapConfig) scrapCFG() error {
	file, err := os.OpenFile("empikCfg.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&s)
	if err != nil {
		panic(err)
	}

	os.Setenv("receiverCSV", s.ReceiverCSV)
	os.Setenv("scrapURL", s.ScrapURL)
	os.Setenv("headerCSV", s.HeaderCSV)

	return nil
}

func (s ScrapConfig) writeHTMLToCSV(name string, price string) error {

	file, err := os.OpenFile(s.ReceiverCSV, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	result := []string{name, price}
	w.Write(result)
	w.Flush()
	return nil
}

func (s *ScrapConfig) scraping() error {
	s.scrapCFG()

	c := colly.NewCollector()

	c.OnHTML("", func(e *colly.HTMLElement) {
		name := e.Attr("")
		fmt.Println(name)
		price := e.Attr("")
		fmt.Println(price)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(s.ScrapURL)
	if err != nil {
		panic(err)
	}
	return nil
}

func main() {
	s := ScrapConfig{}

	file, err := os.OpenFile(s.ReceiverCSV, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	header := strings.Split(s.HeaderCSV, ",")
	headerWriter := csv.NewWriter(file)
	headerWriter.Write(header)
	headerWriter.Flush()
}
