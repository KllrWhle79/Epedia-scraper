package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
)

var CraftsToScrape = [...]string{
	"Armorsmithing",
	"Blacksmithing",
	"Weaponsmithing",
	"Carving",
	"Shaping",
	"Tinkering",
	"Tailoring",
	"Remedies",
}

type CraftItem struct {
	Chapter    int      `json:"chapter"`
	ItemName   string   `json:"itemname"`
	Book       string   `json:"book"`
	Difficulty int      `json:"difficulty"`
	Notes      string   `json:"notes"`
	Technique  string   `json:"technique"`
	Volume     []string `json:"volume"`
}

func ScrapeCraft(craftToScrape string) map[string]CraftItem {
	url := fmt.Sprintf("https://elanthipedia.play.net/%v_products", craftToScrape)
	c := colly.NewCollector()
	// Only used for Weaponsmithing page because the tables are messed up
	chapterNum := 0

	var items = make(map[string]CraftItem)

	// Find and visit all links
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Attr("class") == "wikitable sortable" {
			if strings.Split(e.ChildText("th"), "  ")[0] == "Chapter" ||
				strings.Split(e.ChildText("th"), "  ")[0] == "Item" {
				chapterIndex := -1
				itemIndex := -1
				bookIndex := -1
				difficultyIndex := -1
				notesIndex := -1
				techniqueIndex := -1
				volumeIndex := -1

				e.ForEach("tr", func(i int, f *colly.HTMLElement) {
					if i == 0 {
						chapterNum++
						f.ForEach("th", func(index int, g *colly.HTMLElement) {
							switch strings.TrimSpace(g.Text) {
							case "Chapter":
								chapterIndex = index
							case "Item":
								itemIndex = index
							case "Book":
								bookIndex = index
							case "Difficulty":
								difficultyIndex = index
							case "Notes":
								notesIndex = index
							case "Technique":
								techniqueIndex = index
							case "Volume":
								fallthrough
							case "Vol":
								fallthrough
							case "Ingredients":
								volumeIndex = index
							}
						})
					} else {
						newItem := CraftItem{}
						f.ForEach("td", func(index int, g *colly.HTMLElement) {
							switch index {
							case chapterIndex:
								newItem.Chapter, _ = strconv.Atoi(strings.TrimSpace(g.Text))
							case itemIndex:
								replacer := strings.NewReplacer("<", "", ">", "", "  ", " ")
								newItem.ItemName = replacer.Replace(strings.TrimSpace(g.Text))
							case bookIndex:
								newItem.Book = strings.TrimSpace(g.Text)
							case difficultyIndex:
								reg := regexp.MustCompile(`^(\d+)`)
								newItem.Difficulty, _ = strconv.Atoi(reg.FindString(strings.TrimSpace(g.Text)))
							case notesIndex:
								newItem.Notes = strings.TrimSpace(g.Text)
							case techniqueIndex:
								newItem.Technique = strings.TrimSpace(g.Text)
							case volumeIndex:
								replacer := strings.NewReplacer(" pieces", "",
									" piece", "",
									" and ", "|",
									", ", "|",
									" + ", "|")
								newItem.Volume = strings.Split(replacer.Replace(strings.TrimSpace(g.Text)), "|")
							}
						})
						if chapterIndex == -1 {
							newItem.Chapter = chapterNum
						}
						if newItem.Difficulty != 0 {
							items[newItem.ItemName] = newItem
						}
					}
				})
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(url)
	if err != nil {
		panic(err)
	}

	return items
}
