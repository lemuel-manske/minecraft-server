package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func stopCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the Minecraft service and the EC2 instance (preserves world data)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStop(cmd.Context(), app)
		},
	}
}

func runStop(ctx context.Context, app *App) error {
	inst, err := app.EC2.GetInstance(ctx, app.Config.InstanceName)
	if err != nil {
		return fmt.Errorf("no server found: %w", err)
	}
	if inst.State != "running" {
		return fmt.Errorf("server is not running (state: %s)", inst.State)
	}

	sshClient := app.NewSSH(inst.PublicIP, app.Config.KeyPath)

	s := newSpinner(" Stopping Minecraft service...")
	s.Start()
	if err := sshClient.Run("sudo systemctl stop minecraft"); err != nil {
		s.Stop()
		return fmt.Errorf("stop minecraft: %w", err)
	}
	s.Stop()

	s = newSpinner(" Stopping EC2 instance...")
	s.Start()
	if err := app.EC2.StopInstance(ctx, inst.ID); err != nil {
		s.Stop()
		return fmt.Errorf("stop instance: %w", err)
	}
	s.Stop()

	color.Green("Instance stopped.")
	fmt.Println("World data persists on EBS. Run `mc start` to resume.")
	return nil
}
