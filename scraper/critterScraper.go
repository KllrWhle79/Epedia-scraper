package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

type Critter struct {
	Name      string   `json:"name"`
	MinCap    int      `json:"min_cap"`
	MaxCap    int      `json:"max_cap"`
	BodyType  string   `json:"body_type"`
	Places    []string `json:"places"`
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
				if i > 0 {
					newCritter := Critter{}
					replacer := strings.NewReplacer("<", "", ">", "")
					f.ForEach("td", func(i int, g *colly.HTMLElement) {
						switch g.Attr("class") {
						case "smwtype_wpg":
							newCritter.Name = g.Text
						case "MinCap smwtype_num":
							newCritter.MinCap, _ = strconv.Atoi(replacer.Replace(g.Text))
						case "MaxCap smwtype_num":
							newCritter.MaxCap, _ = strconv.Atoi(replacer.Replace(g.Text))
						case "BodyType smwtype_txt":
							newCritter.BodyType = g.Text
						case "Place smwtype_wpg":
							var places []string
							g.ForEach("a", func(_ int, h *colly.HTMLElement) {
								places = append(places, h.Text)
							})
							newCritter.Places = places
						case "Backstab smwtype_boo":
							newCritter.Backstab, _ = strconv.ParseBool(g.Text)
						case "Gem smwtype_boo":
							newCritter.Gem, _ = strconv.ParseBool(g.Text)
						case "Coin smwtype_boo":
							newCritter.Coin, _ = strconv.ParseBool(g.Text)
						case "Box smwtype_boo":
							newCritter.Box, _ = strconv.ParseBool(g.Text)
						case "Skin smwtype_boo":
							newCritter.Skin, _ = strconv.ParseBool(g.Text)
						case "Cursed smwtype_boo":
							newCritter.Cursed, _ = strconv.ParseBool(g.Text)
						case "Undead smwtype_boo":
							newCritter.Undead, _ = strconv.ParseBool(g.Text)
						case "Construct smwtype_boo":
							newCritter.Construct, _ = strconv.ParseBool(g.Text)
						}
					})
					critters[newCritter.Name] = newCritter
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
