package cmd

import (
	"encoding/json"
	"epediaScraper/scraper"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var Scrape string

func init() {
	runCmd.Flags().StringVarP(&Scrape, "toScrape", "s", "none", "The item to scrape.\n One of the following:\n1. Crafts\n2. Critters")
	RootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		switch Scrape {
		case "crafts":
			var allCraft = make(map[string]map[string]scraper.CraftItem)
			for _, craft := range scraper.CraftsToScrape {
				craftItemsMaps := scraper.ScrapeCraft(craft)
				allCraft[craft] = craftItemsMaps
			}

			jsonData, err := json.MarshalIndent(allCraft, "", "  ")
			if err != nil {
				fmt.Println(err)
			}

			writeJsonToFile("craftItems.json", jsonData)
		case "critters":
			allCritters := scraper.ScrapeCritters()

			jsonData, err := json.MarshalIndent(allCritters, "", "  ")
			if err != nil {
				fmt.Println(err)
			}

			writeJsonToFile("critters.json", jsonData)
		}
	},
}

func writeJsonToFile(fileName string, jsonData []byte) bool {
	filePath := fmt.Sprintf("data/%v", fileName)
	err := ioutil.WriteFile(filePath, jsonData, 0644)
	return err == nil
}
