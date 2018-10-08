package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
)

type Critter struct {
	Name      string   `json:"name"`
	MinCap    int      `json:"min_cap"`
	MaxCap    int      `json:"max_cap"`
	BodyType  string   `json:"body_type"`
	Place     []string `json:"place"`
	Backstab  bool     `json:"backstab"`
	Gem       bool     `json:"gem"`
	Box       bool     `json:"box"`
	Coin      bool     `json:"coin"`
	Skin      bool     `json:"skin"`
	Cursed    bool     `json:"cursed"`
	Undead    bool     `json:"undead"`
	Construct bool     `json:"construct"`
}

func ScrapeCritters() map[string]Critter {
	c := colly.NewCollector()

	var critters = make(map[string]Critter)

	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Attr("class") == "sortable wikitable smwtable" {
			e.ForEach("tr", func(i int, f *colly.HTMLElement) {
				if i == 0 {
					fmt.Println("Header")
				} else {
					fmt.Println("Row")
				}
			})
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit("https://elanthipedia.play.net/Category:Hunting_ladders")
	if err != nil {
		panic(err)
	}

	return critters
}
