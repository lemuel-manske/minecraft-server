package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func backupCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Back up the active server's world",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBackup(cmd.Context(), app)
		},
	}
}

func runBackup(ctx context.Context, app *App) error {
	inst, err := app.EC2.GetInstance(ctx, app.Config.InstanceName)
	if err != nil {
		return fmt.Errorf("no server found: %w", err)
	}
	if inst.State != "running" {
		return fmt.Errorf("server is not running (state: %s)", inst.State)
	}

	profile, err := app.Config.GetProfile(inst.Profile)
	if err != nil {
		return err
	}

	sshClient := app.NewSSH(inst.PublicIP, app.Config.KeyPath)

	s := newSpinner(" Stopping server for backup...")
	s.Start()
	if err := sshClient.Run("sudo systemctl stop minecraft"); err != nil {
		s.Stop()
		return fmt.Errorf("stop minecraft: %w", err)
	}
	s.Stop()

	serverRestarted := false
	defer func() {
		if !serverRestarted {
			rs := newSpinner(" Restarting server after failure...")
			rs.Start()
			_ = sshClient.Run("sudo systemctl start minecraft")
			rs.Stop()
		}
	}()

	s = newSpinner(" Downloading world...")
	s.Start()
	dest := filepath.Join(profile.BackupDir, time.Now().Format("2006-01-02T15-04-05"))
	if err := sshClient.DownloadTar(profile.WorldDir, dest); err != nil {
		s.Stop()
		return fmt.Errorf("download world: %w", err)
	}
	s.Stop()
	color.Green("World downloaded to %s", dest)

	s = newSpinner(" Restarting server...")
	s.Start()
	if err := sshClient.Run("sudo systemctl start minecraft"); err != nil {
		s.Stop()
		return fmt.Errorf("start minecraft: %w", err)
	}
	s.Stop()
	serverRestarted = true

	color.Green("Server restarted.")
	return nil
}
