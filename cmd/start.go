/*
Copyright © 2022 mogenius, Benedikt Iltisberger
*/
package cmd

import (
	"os"
	"os/signal"
	"podloxx-collector/kubernetes"
	"syscall"

	"github.com/mogenius/mo-go/logger"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run the application within your currently selected kubernetes context.",
	Long: `
	Run the application within your currently selected kubernetes context.
	App will cleanup after being terminated with CTRL+C automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		kubernetes.Deploy()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Log.Warning("CLEANUP Kubernetes resources ...")
		kubernetes.Remove()
		logger.Log.Info("CLEANUP finished successfully.")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports start flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
