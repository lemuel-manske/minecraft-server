package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func logsCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "logs",
		Short: "Tail the Minecraft server log",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogs(cmd.Context(), app, os.Stdout)
		},
	}
}

func runLogs(ctx context.Context, app *App, w io.Writer) error {
	inst, err := app.EC2.GetInstance(ctx, app.Config.InstanceName)
	if err != nil {
		return fmt.Errorf("no server found: %w", err)
	}
	if inst.State != "running" {
		return fmt.Errorf("server is not running (state: %s)", inst.State)
	}

	color.Cyan("Streaming logs from %s (Ctrl+C to stop)...\n", inst.PublicIP)

	return app.NewSSH(inst.PublicIP, app.Config.KeyPath).
		Stream("tail -f /home/ubuntu/mc/logs/latest.log", w)
}
