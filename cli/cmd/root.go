package cmd

import (
	"github.com/lkmliz/mc-server/cli/internal/config"
	"github.com/lkmliz/mc-server/cli/internal/ec2"
	"github.com/lkmliz/mc-server/cli/internal/ssh"
	"github.com/lkmliz/mc-server/cli/internal/tf"
	"github.com/spf13/cobra"
)

type App struct {
	Config *config.Config
	EC2    ec2.Client
	TF     tf.Client
	NewSSH func(host, keyPath string) ssh.Client
}

func NewRoot(app *App) *cobra.Command {
	root := &cobra.Command{
		Use:   "mc",
		Short: "Minecraft server manager",
	}

	root.AddCommand(upCmd(app))
	root.AddCommand(downCmd(app))
	root.AddCommand(startCmd(app))
	root.AddCommand(stopCmd(app))
	root.AddCommand(statusCmd(app))
	root.AddCommand(backupCmd(app))
	root.AddCommand(logsCmd(app))

	return root
}
