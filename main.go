package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type ScrapConfig struct {
	ReceiverCSV            string `json:"receiverCSV"`
	ScrapURL               string `json:"scrapURL"`
	HeaderCSV              string `json:"headersCSV"`
	ProductName            string `json:"productName"`
	ProductPrice           string `json:"productPrice"`
	ProductImage           string `json:"productImage"`
	ProductImageAttributes string `json:"productImageAttributes"`
	GoQuery                string `json:"goQuery"`
}

type Scraper interface {
	scrapCFG()
}

func main() {
	if err := doMain(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func doMain() error {
	s := ScrapConfig{}
	path := cfgFileMenu()
	s.scrapCFG(path)

	fmt.Printf("Adding header to %s\n", s.ReceiverCSV)
	header := strings.Split(s.HeaderCSV, ", ")
	s.writeToCSV(header)

	s.scrapHTML()
	fmt.Printf("Succesfully writen to %s\n", s.ReceiverCSV)
	fmt.Println("Exiting program")

	return nil
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
		result = ""
		log.Fatal("Not available yet\nExiting")
	case "3":
		result = ""
		log.Fatal("Not available yet\nExiting")
	}
	return result
}

func searchWord() string {
	var query string
	fmt.Println("Search for product:")
	fmt.Scan(&query)
	return query
}

// environmental variables setup
func (s *ScrapConfig) scrapCFG(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("unable to read cfg file", err, time.Now())
	}

	decoded := map[string]string{}
	err = json.Unmarshal(file, &decoded)
	if err != nil {
		log.Fatalln("impossible to unmarshal cfg", err, time.Now())
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
		case "productName":
			s.ProductName = value
		case "productPrice":
			s.ProductPrice = value
		case "productImage":
			s.ProductImage = value
		case "productImageAttributes":
			s.ProductImageAttributes = value
		case "goQuery":
			s.GoQuery = value
		}
	}
	return nil
}

// writes result to the empikResult.csv file
func (s *ScrapConfig) writeToCSV(result []string) error {
	file, err := os.OpenFile(s.ReceiverCSV, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		log.Fatalln("unable to open receiver file", err, time.Now())
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

	c.OnHTML(s.GoQuery, func(h *colly.HTMLElement) {
		name := h.Attr(s.ProductName)
		price := h.Attr(s.ProductPrice)
		imageURL, _ := h.DOM.Find(s.ProductImage).Attr(s.ProductImageAttributes)
		response := []string{name, price, "zł", imageURL}
		s.writeToCSV(response)
	})

	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		// Visit the next page's URL
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Start scrapping", r.URL.String())
	})

	query := []string{s.ScrapURL, searchWord()}

	err := c.Visit(strings.Join(query, ""))
	if err != nil {
		log.Println("unable to open:", s.ScrapURL, "; error:", err, time.Now())
	}
	return nil
}
