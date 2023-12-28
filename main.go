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
	fmt.Printf("This is a header %s\n", s.HeaderCSV)

	header := strings.Split(s.HeaderCSV, ", ")
	fmt.Println(header)
	s.writeToCSV(header)
	fmt.Println("Succesfully writen")
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

func (s *ScrapConfig) scrapHTML() error {
	c := colly.NewCollector(
		colly.AllowedDomains(s.ScrapURL),
	)

	c.OnHTML("div.search-list-item", func(h *colly.HTMLElement) {
		name := h.Attr("data-product-name")
		price := h.Attr("data-product-price")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	err := c.Visit(s.ScrapURL)
	if err != nil {
		panic("unable to visit URL")
	}
	return nil
}
