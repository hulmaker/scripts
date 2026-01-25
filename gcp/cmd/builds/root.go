package builds

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "builds",
	Short: "Utilities for Cloud Builds.",
}

func init() {
	RootCmd.AddCommand(ListCmd, WatchCmd)
	RootCmd.PersistentFlags().StringVar(&GcloudProject, "project", "", "Google Cloud Project ID (overrides default)")
	RootCmd.PersistentFlags().StringVar(&GcloudRegion, "region", "", "Google Cloud Region (e.g. europe-west1)")
}
