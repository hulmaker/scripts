package cmd

import (
	"fmt"
	"gcp/utils"
	"maps"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

// a safe way to create escaped version of: <base>/<parts/.../>?authuser=1&project=<project_value>
func gcloudUrlBuilder(parts []string) (string, error) {
	return utils.BuildUrl(
		parts, map[string]string{"authuser": "1", "project": project},
	)
}

type urlOption struct {
	url  string
	args []string
}

var baseOptions = map[string]urlOption{
	"console": {"https://console.cloud.google.com/welcome", []string{}},
	"trigger": {"https://console.cloud.google.com/cloud-build/triggers", []string{}},
	"compute": {"https://console.cloud.google.com/compute/instances", []string{}},
	"storage": {"https://console.cloud.google.com/storage/browser", []string{"<path>?"}},
}

// description of keys
var baseOptionsString string = fmt.Sprintf("%v", slices.Collect(maps.Keys(baseOptions)))

// description with args
var baseOptionsDesc string = func(map[string]urlOption) string {
	lines := []string{}
	for key, option := range baseOptions {
		lines = append(lines, fmt.Sprintf("\t%s:\t%v", key, strings.Join(option.args, " ")))
	}
	return strings.Join(lines, "\n")
}(baseOptions)

var openCmd = &cobra.Command{
	Use:   "open " + baseOptionsString + " [args...]",
	Short: "Open the Google Cloud Console in your default browser. See the full help for args description.",
	Long: fmt.Sprintf(`Opens the Google Cloud Console in your default web browser.
Uses the current active project by default, or a specific project/config if specified.

Available args for each option:
%s

Examples:
	gcp open console --project="myproject"
	gcp open compute
	gcp open storage "gcs://bucket/folder" -p "different-project"`, baseOptionsDesc),
	RunE: func(cmd *cobra.Command, args []string) error {
		option, ok := baseOptions[args[0]]
		if !ok {
			return fmt.Errorf("unknown option %q", args[0])
		}
		if len(args)-1 != len(option.args) {
			return fmt.Errorf("too many arguments for option %q - see the full --help", args[0])
		}
		if project != "" && config != "" {
			return fmt.Errorf("can't use --project and --config together")
		}

		if config != "" {
			out, err := utils.CmdGcloudConfigDescribe(config, "get(properties.core.project)").CombinedOutput()
			if err != nil {
				fmt.Printf("Unable to retrieve project for config %s.\nError:\n%s\n", config, string(out))
				fmt.Println("Opening vanilla console.")
			} else {
				project = string(out)
			}
		}

		if project == "" {
			out, err := utils.CmdGcloudGetValue("project").Output()
			if err != nil {
				fmt.Printf("Unable to retrieve current project.\n")
			} else {
				project = string(out)
			}
		}

		// TODO: use switch only for params and then just call OpenBrowser and buildUrl
		args[0] = option.url
		url, err := gcloudUrlBuilder(args)
		if err != nil {
			return err
		}
		return utils.OpenBrowser(url)
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().StringVar(&project, "project", "", "Target Project ID")
	openCmd.Flags().StringVarP(&config, "config", "c", "", "Use Project ID from this gcloud configuration")
}
