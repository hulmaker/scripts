package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

var (
	gcloud_project string
	gcloud_region  string
)

// define color maps
var (
	Reset     = "\033[0m"
	statusMap = map[string]string{
		"SUCCESS":   "\033[1;32m", // Green
		"WORKING":   "\033[33m",   // Yellow
		"QUEUED":    "\033[34m",   // Blue
		"CANCELLED": "\033[35m",   // Magenta
		"FAILURE":   "\033[31m",   // Red
	}
)

func parseAdditionalArgs() []string {
	args := []string{}
	if gcloud_project != "" {
		args = append(args, fmt.Sprintf("--project=%s", gcloud_project))
	}
	if gcloud_region != "" {
		args = append(args, fmt.Sprintf("--region=%s", gcloud_region))
	}
	return args
}

// clamps text to given width
func clamp(text string, width uint) string {
	text = strings.Trim(text, " ")
	w := int(width) // w >= 0
	if len(text) <= w {
		return text + strings.Repeat(" ", w-len(text))
	}
	return text[:w]
}

// Expects: column_name, text. (text is already padded with spaces to the max column width)
type CellFormatter func(string, string) string

// Print the data as a table in human readable format
// Columns are separated by the separator
// cell_formatter to override the default text formatting
// NOTE: text/tabwriter messes up the column alignment when ANSI colors are used.
func printTable(data *[][]string, separator string, formatter CellFormatter) error {
	if len(*data) == 0 {
		return nil
	}

	// calculate the column column_widths
	column_widths := make([]uint, len((*data)[0]))
	for _, row := range *data {
		for j, cell := range row {
			column_widths[j] = max(column_widths[j], uint(len(cell)))
		}
	}

	// print the table
	for _, row := range *data {
		arr := make([]string, len(row))
		for j, text := range row {
			arr[j] = clamp(text, column_widths[j])
			if formatter != nil {
				arr[j] = formatter((*data)[0][j], arr[j])
			}
		}
		fmt.Println(strings.Join(arr, separator))
	}
	return nil
}


func CmdListBuilds(limit uint, filter string, ongoing bool) *exec.Cmd {
	args := []string{
		"builds",
		"list",
		fmt.Sprintf("--limit=%d", limit),
		fmt.Sprintf("--page-size=%d", limit),
		"--format=csv(id,status,substitutions.REPO_NAME,substitutions.BRANCH_NAME,substitutions.TRIGGER_NAME,createTime)",
	}
	if ongoing {
		args = append(args, "--ongoing")
	}
	if filter != "" {
		args = append(args, "--filter="+filter)
	}
	args = append(args, parseAdditionalArgs()...)
	return exec.Command("gcloud", args...)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List and filter Cloud Builds.",
	Long: `List recent cloud builds with formatting options.

Examples:
  gcp builds list -l 10 -w 40 --ongoing --filter="create_time>2026-01-20"

Docs: https://docs.cloud.google.com/build/docs/view-build-results`,
	RunE: func(cmd *cobra.Command, args []string) error {
		width, _ := cmd.Flags().GetUint("column-width")
		limit, _ := cmd.Flags().GetUint("limit")
		filter, _ := cmd.Flags().GetString("filter")
		ongoing, _ := cmd.Flags().GetBool("ongoing")
		out, err := CmdListBuilds(limit, filter, ongoing).CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed requesting builds data with:\n%s", string(out))
		}
		reader := csv.NewReader(strings.NewReader(string(out)))
		records, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("failed parsing CSV: %q", err)
		}
		if len(records) == 0 {
			fmt.Println("no records received")
			return nil
		}
		// crop text to width (except for id column) and color the status column
		formatter := func(column string, text string) string {
			if column == "id" {
				return text
			}

			short_text := clamp(text, min(width, uint(len(text))))
			if column == "status" {
				return statusMap[text] + short_text + Reset
			}
			return short_text
		}
		return printTable(&records, "  ", formatter)
	},
}

type Build struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	FinishTime string `json:"finishTime"`
}

func buildStatus(build_id string, args *[]string) (Build, error) {
	command := exec.Command(
		"gcloud", append([]string{"builds", "describe", build_id, "--format=json"}, (*args)...)...,
	)
	out, err := command.CombinedOutput()
	if err != nil {
		return Build{}, fmt.Errorf("build status request failed with: %q\n%q", err, string(out))
	}
	var build Build
	if err := json.Unmarshal(out, &build); err != nil {
		return Build{}, fmt.Errorf("parsing status response failed with: %q", err)
	}
	return build, nil
}

// Periodically check if the build is done or if it is still going.
func watchBuild(build_id string, max_retries uint, sleep time.Duration, additional_args *[]string) *Build {
	fail_count := uint(0)
	for fail_count < max_retries {
		response, err := buildStatus(build_id, additional_args)
		if err != nil {
			fail_count += 1
			log.Printf("Requesting build status fail count: %d/%d\n", fail_count, max_retries)
			fmt.Println(err)
		} else {
			log.Printf("Build %q status: %q\n", response.ID, response.Status)
			fail_count = 0
		}
		if response.FinishTime != "" {
			return &response
		}
		time.Sleep(sleep * time.Second)
	}
	return nil
}

// run given command and print the stdout to console. If no command is passed, notify.
func runCommand(build *Build, command string) error {
	if command != "" {
		out, err := exec.Command("bash", "-c", command).CombinedOutput()
		fmt.Println(string(out))
		return err
	}

	return beeep.Notify(
		"Build finished",
		fmt.Sprintf("Build %q finished with status: %q", build.ID, build.Status),
		"",
	)
}

var watchCmd = &cobra.Command{
	Use:   "watch <build-id>",
	Short: "Monitor a build and trigger actions on completion.",
	Long:  `Polls the status of a specific build. Can send a desktop notification or run a shell command when the build finishes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("please provide the build <id>")
		}
		build_id := args[0]
		sleep, _ := cmd.Flags().GetUint("sleep")
		max_retries, _ := cmd.Flags().GetUint("retry")
		command_string, _ := cmd.Flags().GetString("command")
		only_success, _ := cmd.Flags().GetBool("only-success")
		additional_args := parseAdditionalArgs()

		build := watchBuild(build_id, max_retries, time.Duration(sleep), &additional_args)
		if build == nil {
			return fmt.Errorf("failed requesting the build status")
		}

		if only_success && build.Status != "SUCCESS" {
			fmt.Printf("Build did not finish with SUCCESS (--only-success=%t) exiting...\n", only_success)
			return nil
		}
		return runCommand(build, command_string)
	},
}

var buildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "Utilities for Cloud Builds.",
}

func init() {
	rootCmd.AddCommand(buildsCmd)
	buildsCmd.AddCommand(listCmd, watchCmd)
	buildsCmd.PersistentFlags().StringVar(&gcloud_project, "project", "", "Google Cloud Project ID (overrides default)")
	buildsCmd.PersistentFlags().StringVar(&gcloud_region, "region", "", "Google Cloud Region (e.g. europe-west1)")

	listCmd.Flags().UintP("column-width", "w", 40, "Max width for output columns")
	listCmd.Flags().UintP("limit", "l", 50, "Max number of builds to list")
	listCmd.Flags().BoolP("ongoing", "o", false, "Show only unfinished builds")
	listCmd.Flags().String("filter", "", "Filter expression (e.g. 'status=SUCCESS')")

	watchCmd.Flags().UintP("sleep", "s", 60, "Polling interval in seconds")
	watchCmd.Flags().String("command", "", "Shell command to run upon completion")
	watchCmd.Flags().BoolP("only-success", "o", false, "Only run command if build succeeds")
	watchCmd.Flags().UintP("retry", "r", 5, "Max retries for status checks")
}
