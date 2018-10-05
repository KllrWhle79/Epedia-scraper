package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use: "epedia-scraper",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
