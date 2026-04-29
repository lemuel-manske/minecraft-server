package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func startCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start a stopped EC2 instance and resume the Minecraft service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(cmd.Context(), app)
		},
	}
}

func runStart(ctx context.Context, app *App) error {
	inst, err := app.EC2.GetInstance(ctx, app.Config.InstanceName)
	if err != nil {
		return fmt.Errorf("no server found — run `mc up` first: %w", err)
	}

	s := newSpinner(" Starting EC2 instance...")
	s.Start()
	if err := app.EC2.StartInstance(ctx, inst.ID); err != nil {
		s.Stop()
		return fmt.Errorf("start instance: %w", err)
	}

	ip, err := app.EC2.WaitRunning(ctx, inst.ID)
	if err != nil {
		s.Stop()
		return fmt.Errorf("wait for instance: %w", err)
	}
	s.Stop()

	sshClient := app.NewSSH(ip, app.Config.KeyPath)

	s = newSpinner(" Waiting for SSH...")
	s.Start()
	if err := sshClient.WaitForSSH(ctx); err != nil {
		s.Stop()
		return fmt.Errorf("wait for SSH: %w", err)
	}
	s.Stop()

	s = newSpinner(" Starting Minecraft service...")
	s.Start()
	if err := sshClient.Run("sudo systemctl start minecraft"); err != nil {
		s.Stop()
		return fmt.Errorf("start minecraft: %w", err)
	}
	s.Stop()

	color.Green("Server is up at %s:25565", ip)
	return nil
}
