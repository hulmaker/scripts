package cmd

import (
	"fmt"
	"strings"
	"net/url"
	"github.com/spf13/cobra"
	"gcp/utils"
)

var (
	// TODO: Is it global as in the whole project? If I name it simple project, I get duplicate warning
	userProject string
	configName string
)

// a safe way to create escaped version of: <base>?project=<project_value>&authuser=1
func buildConsoleUrl(project string) string {
	// https://console.cloud.google.com/welcome?project=named-purpose-220220&authuser=1
	base := "https://console.cloud.google.com/welcome"
	if project == "" {
		return base
	}
	u, _:= url.Parse(base)
	query := u.Query()
	query.Set("authuser", "1")
	query.Set("project", strings.TrimSpace(project))
	u.RawQuery = query.Encode()
	return u.String()
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Open the Google Cloud Console in browser.",
	Long: `Opens the Google Cloud Console in your default web browser.
Targets the current active project by default, or a specific project/config if specified.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if userProject != "" && configName != "" {
			return fmt.Errorf("can't use --project and --config together")
		}

		if configName != "" {
			out, err := utils.CmdGcloudConfigDescribe(configName, "get(properties.core.project)").CombinedOutput()
			if err != nil {
				fmt.Printf("Unable to retrieve project for config %s.\nError:\n%s\n", configName, string(out))
				fmt.Println("Opening vanilla console.")
			} else {
				userProject = string(out)
			}
		}

		if userProject == "" {
			out, err := utils.CmdGcloudGetValue("project").Output()
			if err != nil {
				fmt.Printf("Unable to retrieve current project, opening vanilla console.\n")
			} else {
				userProject = string(out)
			}
		}
		return utils.OpenBrowser(buildConsoleUrl(userProject))
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)
	consoleCmd.Flags().StringVar(&userProject, "project", "", "Target Project ID")
	consoleCmd.Flags().StringVarP(&configName, "config", "c", "", "Use Project ID from this gcloud configuration")
}
