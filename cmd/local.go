/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"podloxx-collector/api"
	"podloxx-collector/network"

	"github.com/spf13/cobra"
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Run the application on your local machine (root privileges required).",
	Long: `
	Run the application on your local machine's network devices. 
	Awesome hacker view window. Should always be visible to impress non-it-folks.`,
	Run: func(cmd *cobra.Command, args []string) {
		go network.MonitorAll(true, "")
		api.InitApi()

		// quit := make(chan os.Signal)
		// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		// <-quit
		// logger.Log.Warning("XXX CLEANUP HERE")
	},
}

func init() {
	rootCmd.AddCommand(localCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
