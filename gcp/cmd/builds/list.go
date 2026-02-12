package builds

import (
	"encoding/csv"
	"fmt"
	"gcp/utils"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var reset = "\033[0m"
var statusMap = map[string]string{
	"SUCCESS":   "\033[1;32m", // Green
	"WORKING":   "\033[33m",   // Yellow
	"QUEUED":    "\033[34m",   // Blue
	"CANCELLED": "\033[35m",   // Magenta
	"FAILURE":   "\033[31m",   // Red
}

func cmdListBuilds(limit uint, filter string, ongoing bool) *exec.Cmd {
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
	args = append(args, ParseAdditionalArgs()...)
	return exec.Command("gcloud", args...)
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List and filter Cloud Builds.",
	Long: `List recent cloud builds with formatting options.

Examples:
  gcp builds list -l 10 -w 40 --ongoing --filter="create_time>2026-01-20"

Note: Using filter is reccommended with --list. Filter makes the fetch significantly slower.
Docs: https://docs.cloud.google.com/build/docs/view-build-results`,
	RunE: func(cmd *cobra.Command, args []string) error {
		width, _ := cmd.Flags().GetUint("column-width")
		limit, _ := cmd.Flags().GetUint("limit")
		filter, _ := cmd.Flags().GetString("filter")
		ongoing, _ := cmd.Flags().GetBool("ongoing")
		out, err := cmdListBuilds(limit, filter, ongoing).CombinedOutput()
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

			shortText := utils.Clamp(text, min(width, uint(len(text))))
			if column == "status" {
				return statusMap[text] + shortText + reset
			}
			return shortText
		}
		return utils.PrintTable(records, "  ", formatter)
	},
}

func init() {
	ListCmd.Flags().UintP("column-width", "w", 40, "Max width for output columns")
	ListCmd.Flags().UintP("limit", "l", 50, "Max number of builds to list")
	ListCmd.Flags().BoolP("ongoing", "o", false, "Show only unfinished builds")
	ListCmd.Flags().String("filter", "", "Filter expression (e.g. 'status=SUCCESS')")
}
