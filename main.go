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
	SearchWord  string `json:"searchWord"`
}

type Scraper interface {
	scrapCFG() string
}

func main() {
	s := ScrapConfig{}
	path := "empikCfg.json"
	s.scrapCFG(path)

	header := strings.Split(s.HeaderCSV, ", ")
	s.writeToCSV(header)

	s.scrapHTML()
	fmt.Printf("Succesfully writen to %s\n", s.ReceiverCSV)
	fmt.Println("Exiting program")
}

func (s *ScrapConfig) scrapCFG(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		panic("unable to read cfg file")
	}

	decoded := map[string]string{}
	err = json.Unmarshal(file, &decoded)
	if err != nil {
		panic("unmarshall is not possible")
	}

	for option, value := range decoded {
		os.Setenv(strings.ToUpper(option), value)
		switch option {
		case "headersCSV":
			s.HeaderCSV = value
		case "scrapURL":
			s.ScrapURL = value
		case "receiverCSV":
			s.ReceiverCSV = value
		case "searchWord":
			s.SearchWord = value
		}
	}
	return nil
}

func (s *ScrapConfig) writeToCSV(result []string) error {
	file, err := os.OpenFile(s.ReceiverCSV, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic("unable to open receiverCSV")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write(result)
	writer.Flush()
	return nil
}

// scrap required URL and write result into CSV file
func (s *ScrapConfig) scrapHTML() error {
	c := colly.NewCollector()

	c.OnHTML("div.search-list-item", func(h *colly.HTMLElement) {
		name := h.Attr("data-product-name")
		price := h.Attr("data-product-price")
		imageURL, _ := h.DOM.Find("meta[itemprop=\"image\"]").Attr("content")
		response := []string{name, price, "zł", imageURL}
		s.writeToCSV(response)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Start scrapping", r.URL.String())
	})

	err := c.Visit(s.ScrapURL)
	if err != nil {
		panic("unable to visit URL")
	}
	return nil
}
