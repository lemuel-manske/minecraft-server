package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func statusCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show server and instance status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(cmd.Context(), app)
		},
	}
}

func runStatus(ctx context.Context, app *App) error {
	inst, err := app.EC2.GetInstance(ctx, app.Config.InstanceName)
	if err != nil {
		color.Yellow("No server running.")
		return nil
	}

	bold := color.New(color.Bold).SprintFunc()
	fmt.Printf("%s  %s\n", bold("profile:"), inst.Profile)
	fmt.Printf("%s       %s\n", bold("ip:"), inst.PublicIP)
	fmt.Printf("%s      %s (%s)\n", bold("ec2:"), inst.ID, inst.State)

	if inst.State != "running" {
		return nil
	}

	sshClient := app.NewSSH(inst.PublicIP, app.Config.KeyPath)
	out, err := sshClient.Output("sudo systemctl status minecraft")
	if err != nil {
		return fmt.Errorf("service status: %w", err)
	}

	fmt.Println()
	fmt.Println(out)
	return nil
}
