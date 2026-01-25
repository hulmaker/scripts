package builds

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

type Build struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	FinishTime string `json:"finishTime"`
}

var (
	GcloudProject string
	GcloudRegion  string
)

func ParseAdditionalArgs() []string {
	args := []string{}
	if GcloudProject != "" {
		args = append(args, fmt.Sprintf("--project=%s", GcloudProject))
	}
	if GcloudRegion != "" {
		args = append(args, fmt.Sprintf("--region=%s", GcloudRegion))
	}
	return args
}


func buildStatus(buildID string, args *[]string) (Build, error) {
	command := exec.Command(
		"gcloud", append([]string{"builds", "describe", buildID, "--format=json"}, (*args)...)...,
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
func watchBuild(buildID string, maxRetries uint, sleep time.Duration, additionalArgs *[]string) *Build {
	failCount := uint(0)
	for failCount < maxRetries {
		response, err := buildStatus(buildID, additionalArgs)
		if err != nil {
			failCount += 1
			log.Printf("Requesting build status fail count: %d/%d\n", failCount, maxRetries)
			fmt.Println(err)
		} else {
			log.Printf("Build %q status: %q\n", response.ID, response.Status)
			failCount = 0
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
		out, err := exec.Command("sh", "-c", command).CombinedOutput()
		fmt.Println(string(out))
		return err
	}

	return beeep.Notify(
		"Build finished",
		fmt.Sprintf("Build %q finished with status: %q", build.ID, build.Status),
		"",
	)
}

var WatchCmd = &cobra.Command{
	Use:   "watch <build-id>",
	Short: "Monitor a build and trigger actions on completion.",
	Long:  `Polls the status of a specific build. Can send a desktop notification or run a shell command when the build finishes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("please provide the build <id>")
		}
		buildID := args[0]
		sleep, _ := cmd.Flags().GetUint("sleep")
		maxRetries, _ := cmd.Flags().GetUint("retry")
		commandString, _ := cmd.Flags().GetString("command")
		onlySuccess, _ := cmd.Flags().GetBool("only-success")
		additionalArgs := ParseAdditionalArgs()

		build := watchBuild(buildID, maxRetries, time.Duration(sleep), &additionalArgs)
		if build == nil {
			return fmt.Errorf("failed requesting the build status")
		}

		if onlySuccess && build.Status != "SUCCESS" {
			fmt.Printf("Build did not finish with SUCCESS (--only-success=%t) exiting...\n", onlySuccess)
			return nil
		}
		return runCommand(build, commandString)
	},
}

func init() {
	WatchCmd.Flags().UintP("sleep", "s", 60, "Polling interval in seconds")
	WatchCmd.Flags().String("command", "", "Shell command to run upon completion")
	WatchCmd.Flags().BoolP("only-success", "o", false, "Only run command if build succeeds")
	WatchCmd.Flags().UintP("retry", "r", 5, "Max retries for status checks")
}
