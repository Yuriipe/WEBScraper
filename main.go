package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

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

// setting env variables
func scrapCFG(path string) error {
	s := &ScrapConfig{}
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
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

// writing scrap result to csv
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

// scraping from HTML
func (s *ScrapConfig) scraping() error {
	scrapCFG()

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
	s.scrapCFG("empikCfg.csv")

	file, err := os.OpenFile(s.ReceiverCSV, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	w := csv.NewWriter(file)
	tekst := []string{"new tekst", "to csv"}
	w.Write(tekst)
	w.Flush()

	// file, err := os.Open(s.ReceiverCSV)
	// if err != nil {
	// 	panic(err)
	// }

	// header := strings.Split(s.HeaderCSV, ",")
	// headerWriter := csv.NewWriter(file)
	// headerWriter.Write(header)
	// headerWriter.Flush()
}
