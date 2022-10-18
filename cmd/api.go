/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"podloxx-collector/api"

	"github.com/mogenius/mo-go/utils"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run the application in API mode.",
	Long: `
	In API mode you can use all gathered data from the websocket api and REST api.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.OpenBrowser("http://127.0.0.1:8080/traffic")
		api.InitApi()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
