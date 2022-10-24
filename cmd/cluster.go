/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"podloxx-collector/api"
	"podloxx-collector/network"

	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Run the application on your cluster (root privileges required).",
	Long: `
	Run the application on your cluster machine's network devices. 
	Awesome hacker view window. Should always be visible to impress non-it-folks.`,
	Run: func(cmd *cobra.Command, args []string) {
		go network.MonitorAll(false, "")
		api.InitApiCluster()
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports cluster flags which will only run when this command
	// is called directly, e.g.:
	// clusterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
