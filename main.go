package main

import (
	"encoding/csv"
	"os"

	"golang.org/x/text/collate"
)

type ScrapConfig struct {
	URL   string
	Query string
}

type Marketplace struct {
	Allegro []ScrapConfig
	Olx     []ScrapConfig
	Amazon  []ScrapConfig
}

type WEBScrapper struct {
	cfg         ScrapConfig
	marketplace Marketplace
}

func (w WEBScrapper) scraping(m Marketplace) []string {
	c := colly.NewCollector()

}

func main() {
	if err := doMain(); err != nil {
		panic(err)
	}

}

func doMain() {

	file, err := os.Create("result.csv")
	if err != nil {
		panic("unable to create result.csv")
	}
	defer file.Close()

	writer := csv.NewWriter(file)

}
