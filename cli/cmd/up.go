package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func upCmd(app *App) *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:   "up <profile>",
		Short: "Provision and start the server with the given profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUp(cmd.Context(), app, args[0], verbose)
		},
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show full terraform output")
	return cmd
}

func runUp(ctx context.Context, app *App, profileName string, verbose bool) error {
	profile, err := app.Config.GetProfile(profileName)
	if err != nil {
		return err
	}

	if _, err := os.Stat(profile.ModpackDir); err != nil {
		return fmt.Errorf("modpack dir %q not found: %w", profile.ModpackDir, err)
	}

	if verbose {
		app.TF.SetVerbose(true)
	}

	vars := buildTFVars(app, profileName, profile.ForgeURL, profile.ForgeSHA256)

	s := newSpinner(" Provisioning infrastructure...")
	if !verbose {
		s.Start()
	}
	if err := app.TF.Apply(ctx, vars); err != nil {
		s.Stop()
		return fmt.Errorf("terraform apply: %w", err)
	}
	s.Stop()

	ip, err := app.TF.Output(ctx, "instance_ip")
	if err != nil {
		return fmt.Errorf("get instance IP: %w", err)
	}

	sshClient := app.NewSSH(ip, app.Config.KeyPath)

	s = newSpinner(" Waiting for SSH...")
	s.Start()
	if err := sshClient.WaitForSSH(ctx); err != nil {
		s.Stop()
		return fmt.Errorf("wait for SSH: %w", err)
	}
	s.Stop()

	s = newSpinner(" Uploading modpack...")
	s.Start()
	if err := sshClient.UploadTar(profile.ModpackDir, "/home/ubuntu/mc"); err != nil {
		s.Stop()
		return fmt.Errorf("upload modpack: %w", err)
	}
	s.Stop()
	color.Green("Modpack uploaded.")

	props := filepath.Join(profile.ModpackDir, "server.properties")
	s = newSpinner(" Uploading server config...")
	s.Start()
	if err := sshClient.UploadFile(props, "/home/ubuntu/mc/server.properties"); err != nil {
		s.Stop()
		return fmt.Errorf("upload server.properties: %w", err)
	}
	s.Stop()
	color.Green("Server config uploaded.")

	if err := sshClient.Run("sudo systemctl start minecraft"); err != nil {
		return fmt.Errorf("start minecraft: %w", err)
	}

	color.Green("Server is up at %s:25565", ip)
	return nil
}

func newSpinner(suffix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = suffix
	return s
}

func buildTFVars(app *App, profileName, forgeURL, forgeSHA256 string) []string {
	quoted := make([]string, len(app.Config.AllowedCIDRs))
	for i, c := range app.Config.AllowedCIDRs {
		quoted[i] = `"` + c + `"`
	}

	return []string{
		"operator_cidr=" + app.Config.OperatorCIDR,
		"allowed_cidrs=[" + strings.Join(quoted, ",") + "]",
		"forge_url=" + forgeURL,
		"forge_sha256=" + forgeSHA256,
		"profile_name=" + profileName,
	}
}
