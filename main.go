package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Product struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	ImgURL string `json:"imgurl"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.empik.com"),
	)

	c.OnHTML("div.search-list-item", func(h *colly.HTMLElement) {
		fmt.Println(h.Attr("data-product-name"))
		fmt.Println(h.Attr("data-product-price"), "zł")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	err := c.Visit("https://www.empik.com/szukaj/produkt?q=slime&qtype=basicForm&trending=true")
	if err != nil {
		panic(err)
	}
}
