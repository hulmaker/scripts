package builds

import "github.com/spf13/cobra"

var (
	project string
	region  string
)

func ParseAdditionalArgs() []string {
	args := []string{}
	if project != "" {
		args = append(args, "--project="+project)
	}
	if region != "" {
		args = append(args, "--region="+region)
	}
	return args
}

var RootCmd = &cobra.Command{
	Use:   "builds",
	Short: "Utilities for listing and actively watching cloud builds.",
}

func init() {
	RootCmd.AddCommand(ListCmd, WatchCmd)
	RootCmd.PersistentFlags().StringVar(&project, "project", "", "Google Cloud Project ID (overrides default)")
	RootCmd.PersistentFlags().StringVar(&region, "region", "", "Google Cloud Region (e.g. europe-west1)")
}
