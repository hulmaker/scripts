/*
Here we store various gcloud commands that are likely to be reused.
*/
package utils

import (
	"io"
	"os/exec"
)


// get value like current "project" or "account"
func CmdGcloudGetValue(value string) *exec.Cmd {
	cmd := exec.Command("gcloud", "config", "get-value", "project")
	cmd.Stderr = io.Discard // this command often prints nonsense to output.
	return cmd
}


// list existing configurations
func CmdGcloudConfigList() *exec.Cmd {
	return exec.Command("gcloud", "config", "configurations", "list")
}


// info about existing config. (format example: "get(properties.core.project)")
func CmdGcloudConfigDescribe(configName string, format string) *exec.Cmd {
	args := []string{"config", "configurations", "describe", configName}
	if format != "" {
		args = append(args, "--format=" + format)
	}
	return exec.Command("gcloud", args...)
}
