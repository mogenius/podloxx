/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"os"
	"os/signal"
	"podloxx-collector/api"
	"syscall"

	"github.com/mogenius/mo-go/logger"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run the application within your currently selected kubernetes context.",
	Long: `
	Run the application within your currently selected kubernetes context.
	App will cleanup after being terminated with CTRL+C automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		api.TestRedis()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Log.Info("CLEANUP finished successfully.")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports test flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
