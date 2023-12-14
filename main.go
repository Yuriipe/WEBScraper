package main

import (
	"github.com/gocolly/colly"
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

type Product struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	ImgURL string `json:"imgurl"`
}

func (w WEBScrapper) scraping(m Marketplace) []string {
	c := colly.NewCollector(
		colly.AllowedDomains("allegro.pl", "amazon.pl", "olx.pl"),
	)

	c.OnHTML("", func(e *colly.HTMLElement) {

	})
}

func main() {
	if err := doMain(); err != nil {
		panic(err)
	}

}

func doMain() {

}
