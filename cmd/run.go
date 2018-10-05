package cmd

import (
	"encoding/json"
	"epediaScraper/scraper"
	"fmt"
	"github.com/spf13/cobra"
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

var Scrape string

var craftMap struct {
	name string
	craftItem
}

func init() {
	runCmd.Flags().StringVarP(&Scrape, "toScrape", "s", "none", "The item to scrape.\n One of the following:\n1. Crafts\n2. Critters")
	RootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		switch Scrape {
		case "Crafts":
			var allCraft = map[string]interface{}{}
			for _, craft := range CraftsToScrape {
				craftItemsMaps := scraper.ScrapeCraft(craft)
				allCraft[craft] = craftItemsMaps
			}

			test, err := json.Marshal(allCraft)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(test)
			}
		case "Critters":

		}
	},
}
