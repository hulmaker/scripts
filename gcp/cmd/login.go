package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"gcp/utils"
	"github.com/spf13/cobra"
)

func cmdConfigActivate() *exec.Cmd {
	return exec.Command("gcloud", "config", "configurations", "activate", config)
}

func cmdConfigCreate() *exec.Cmd {
	return exec.Command("gcloud", "config", "configurations", "create", config)
}

func validateAccountProjectFlags() error {
	if project != "" && strings.Contains(project, "@") {
		return fmt.Errorf("--project %q looks like an email; did you mean --account?", project)
	}
	if account != "" && !strings.Contains(account, "@") {
		return fmt.Errorf("--account %q does not look like an email; did you mean --project?", account)
	}
	return nil
}

func activateCreateActivate() error {
	out, err := cmdConfigActivate().CombinedOutput()
	if err == nil {
		return nil
	}
	if !strings.Contains(string(out), "does not exist") {
		return fmt.Errorf("activating config %s failed with:\n%s", config, string(out))
	}

	if project == "" || account == "" {
		out, err = utils.CmdGcloudConfigList().Output()
		msg := "config %q not found\nto create a new config, provide --account and --project"
		if err != nil {
			return fmt.Errorf(msg, config)
		}
		return fmt.Errorf(msg+" (or see the available configs below)\n\n%s", config, string(out))
	}

	out, err = cmdConfigCreate().CombinedOutput()
	if err != nil {
		return fmt.Errorf("creating config %s failed with:\n%s", config, string(out))
	}

	fmt.Printf("created new config: %s project: %s account %s\n", config, project, account)
	out, err = cmdConfigActivate().CombinedOutput()
	if err != nil {
		return fmt.Errorf("activating config %s failed with:\n%s", config, string(out))
	}
	return nil
}

// TODO: gcp switch -> switch between existing configurations
//
//	login() set's project, account and logs in with default credentials
//	this creates application_default_credentials.json that we need for running whatever
//	But we need to open browser every time.
//	It is possible to define configurations and switch between them (see the code below)
//	The problem is they don't switch default credentials.
//	To solve this, login script must do the following:
//	    1. create new configuration with given name if it does not exist
//	    2. activate the configuration
//	    3. copy application_default_credentials.json" .gcloud-configs/${config}-adc.json"
//	Switching configurations then does the following:
//	    1. switch to existing config that was created during login.
//	    2. either set GOOGLE_APPLICATION_CREDENTIALS to fuul path to ${config}-adc.json
//	    3. or set application_default_credentials.json to contents in ${config}-adc.json
//	    4. setting only envvar does not work for applications running in docker.
//	Below is some helpful code:
//
//	config_dir="${HOME}/.configs/.gcloud-configs"
//	config="rs" # do this for each config
//	gcloud config configurations create "$config"
//	gcp login "$config" (existing function in this script)
//	cp "${HOME}/.config/gcloud/application_default_credentials.json" "${config_dir}/${config}-adc.json"
func login(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("specify the config name you want to activate")
		// todo if is active, do nothing
	}
	if err := validateAccountProjectFlags(); err != nil {
		return err
	}
	config = args[0]
	err := activateCreateActivate()
	if err != nil {
		return err
	}

	commands := [][]string{}
	if account != "" {
		commands = append(commands, []string{"gcloud", "config", "set", "account", account})
	}
	if project != "" {
		commands = append(commands, []string{"gcloud", "config", "set", "project", project})
	}
	adcLogin := []string{"gcloud", "auth", "application-default", "login"}
	if account != "" {
		adcLogin = append(adcLogin, "--account="+account)
	}
	commands = append(commands,
		adcLogin,
		[]string{"gcloud", "config", "set", "disable_usage_reporting", "true"},
	)

	err = utils.RunCommands(commands)
	if err != nil {
		return fmt.Errorf("failed setting application-default")
	}
	return nil
}

var loginCmd = &cobra.Command{
	Use:   "login <config-name>",
	Short: "Log in to a gcloud SDK and set default application login.",
	Long: `Authenticate with Google Cloud and setup a specific configuration with project and account settings.
It ensures application default credentials are set correctly for local development.

List existing configurations with: gcloud config configurations list`,
	RunE: login,
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVar(&project, "project", "", "Google Cloud Project ID")
	loginCmd.Flags().StringVar(&account, "account", "", "Google Cloud Account Email")
}
