package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	fmt.Printf("this is %s", s.HeaderCSV)

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
		fmt.Printf("Variable %s was set to value: %s\n", option, value)
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
