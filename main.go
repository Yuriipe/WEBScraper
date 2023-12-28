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
}

// add user input for several options:
// 3. cronjob possibility
// 4. output file type choise

type Scraper interface {
	scrapCFG()
}

func main() {
	s := ScrapConfig{}
	path := cfgFileMenu()
	s.scrapCFG(path)

	header := strings.Split(s.HeaderCSV, ", ")
	s.writeToCSV(header)

	s.scrapHTML()
	fmt.Printf("Succesfully writen to %s\n", s.ReceiverCSV)
	fmt.Println("Exiting program")
}

// user chooses the URL to scrap
func cfgFileMenu() string {
	var choice string
	fmt.Println("Choose required URL:")
	fmt.Println("1. empik")
	fmt.Println("2. amazon")
	fmt.Println("3. allegro")
	fmt.Scanln(&choice)

	var result string
	switch choice {
	case "1":
		result = "empikCfg.json"
		fmt.Println("Running empik scrap")
	case "2":
		result = "not available yet"
		fmt.Println("Not available yet\nExiting")
		os.Exit(0)
	case "3":
		result = "not available yet"
		fmt.Println("Not available yet\nExiting")
		os.Exit(0)
	}
	return result
}

func searchWord() string {
	var query string
	fmt.Println("What product would you like to search for?\n(for search improvement give one word description)")
	fmt.Scanln(&query)
	return query
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
		}
	}
	return nil
}

// writes result to the empikResult.csv file
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

// scraps required URL
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

	query := []string{s.ScrapURL, searchWord()}

	err := c.Visit(strings.Join(query, ""))
	if err != nil {
		panic("unable to visit URL")
	}
	return nil
}
