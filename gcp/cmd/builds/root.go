package builds

import "github.com/spf13/cobra"


var (
	GcloudProject string
	GcloudRegion  string
)

func ParseAdditionalArgs() []string {
	args := []string{}
	if GcloudProject != "" {
		args = append(args, "--project="+GcloudProject)
	}
	if GcloudRegion != "" {
		args = append(args, "--region="+GcloudRegion)
	}
	return args
}


var RootCmd = &cobra.Command{
	Use:   "builds",
	Short: "Utilities for listing and actively watching cloud builds.",
}

func init() {
	RootCmd.AddCommand(ListCmd, WatchCmd)
	RootCmd.PersistentFlags().StringVar(&GcloudProject, "project", "", "Google Cloud Project ID (overrides default)")
	RootCmd.PersistentFlags().StringVar(&GcloudRegion, "region", "", "Google Cloud Region (e.g. europe-west1)")
}
