package cmd

import (
	"fmt"
	"strings"
	"io"
	"net/url"
	"os/exec"
	"runtime"
	"github.com/spf13/cobra"
)

var (
	// TODO: Is it global as in the whole project? If I name it simple project, I get duplicate warning
	user_project string
	config_name string
)

func openBrowser(url string) error {
	fmt.Printf("Opening url in default browser. If nothing happens, plase visit:\n%s\n", url)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

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
	// query.Set("penisoid", "jahodnik")
	u.RawQuery = query.Encode()
	return u.String()
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Open the Google Cloud Console in browser.",
	Long: `Opens the Google Cloud Console in your default web browser.
Targets the current active project by default, or a specific project/config if specified.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if user_project != "" && config_name != "" {
			return fmt.Errorf("can't use --project and --config together")
		}

		if config_name != "" {
			cmd := exec.Command("gcloud", "config", "configurations", "describe", config_name, "--format=get(properties.core.project)")
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Unable to retrieve project for config %s.\nError:\n%s\n", config_name, string(out))
				fmt.Println("Opening vanilla console.")
			} else {
				user_project = string(out)
			}
		}

		if user_project == "" {
			cmd := exec.Command("gcloud", "config", "get-value", "project")
			cmd.Stderr = io.Discard
			out, err := cmd.Output()
			if err != nil {
				fmt.Printf("Unable to retrieve current project, opening vanilla console.\n")
			} else {
				user_project = string(out)
			}
		}
		return openBrowser(buildConsoleUrl(user_project))
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)
	consoleCmd.Flags().StringVar(&user_project, "project", "", "Target Project ID")
	consoleCmd.Flags().StringVarP(&config_name, "config", "c", "", "Use Project ID from this gcloud configuration")
}
