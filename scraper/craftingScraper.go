package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
)

type CraftItem struct {
	chapter    int
	item       string
	book       string
	difficulty int
	notes      string
	technique  string
	volume     []string
}

func ScrapeCraft(craftToScrape string) map[string]*CraftItem {
	url := fmt.Sprintf("https://elanthipedia.play.net/%v_products", craftToScrape)
	c := colly.NewCollector()
	firstColumn := ""
	// Only used for Weaponsmithing page because the tables are messed up
	chapterNum := 0

	if craftToScrape == "Weaponsmithing" {
		firstColumn = "Item"
	} else {
		firstColumn = "Chapter"
	}

	var items = map[string]*CraftItem{}

	// Find and visit all links
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Attr("class") == "wikitable sortable" {
			if strings.Split(e.ChildText("th"), "  ")[0] == firstColumn {
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
						newItem := &CraftItem{}
						f.ForEach("td", func(index int, g *colly.HTMLElement) {
							switch index {
							case chapterIndex:
								newItem.chapter, _ = strconv.Atoi(strings.TrimSpace(g.Text))
							case itemIndex:
								replacer := strings.NewReplacer("<", "", ">", "", "  ", " ")
								newItem.item = replacer.Replace(strings.TrimSpace(g.Text))
							case bookIndex:
								newItem.book = strings.TrimSpace(g.Text)
							case difficultyIndex:
								reg := regexp.MustCompile(`^(\d+)`)
								newItem.difficulty, _ = strconv.Atoi(reg.FindString(strings.TrimSpace(g.Text)))
							case notesIndex:
								newItem.notes = strings.TrimSpace(g.Text)
							case techniqueIndex:
								newItem.technique = strings.TrimSpace(g.Text)
							case volumeIndex:
								newItem.volume = strings.Split(strings.TrimSpace(g.Text), ", ")
							}
						})
						if chapterIndex == -1 {
							newItem.chapter = chapterNum
						}
						if newItem.difficulty != 0 {
							items[newItem.item] = newItem
						}
					}
				})
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)

	return items
}
