package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func downCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "down <profile>",
		Args:  cobra.ExactArgs(1),
		Short: "Stop the server and destroy infrastructure",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDown(cmd.Context(), app, args[0])
		},
	}
}

func runDown(ctx context.Context, app *App, profileName string) error {
	profile, err := app.Config.GetProfile(profileName)
	if err != nil {
		return err
	}
	vars := buildTFVars(app, profileName, profile)

	inst, err := app.EC2.GetInstance(ctx, app.Config.InstanceName)
	if err != nil {
		return fmt.Errorf("no server found: %w", err)
	}

	if inst.State == "running" {
		sshClient := app.NewSSH(inst.PublicIP, app.Config.KeyPath)
		s := newSpinner(" Stopping Minecraft service...")
		s.Start()
		_ = sshClient.Run("sudo systemctl stop minecraft")
		s.Stop()
	}

	s := newSpinner(" Destroying infrastructure...")
	s.Start()
	if err := app.TF.Destroy(ctx, vars); err != nil {
		s.Stop()
		return fmt.Errorf("terraform destroy: %w", err)
	}
	s.Stop()

	color.Green("Infrastructure destroyed.")
	return nil
}
