package main

import (
	"fmt"
	"os"

	"github.com/lkmliz/mc-server/cli/cmd"
	"github.com/lkmliz/mc-server/cli/internal/config"
	"github.com/lkmliz/mc-server/cli/internal/ec2"
	"github.com/lkmliz/mc-server/cli/internal/ssh"
	"github.com/lkmliz/mc-server/cli/internal/tf"
)

func main() {
	cfg, err := config.Load("mc.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "load config:", err)
		os.Exit(1)
	}

	ec2Client, err := ec2.New(cfg.Region)
	if err != nil {
		fmt.Fprintln(os.Stderr, "init ec2:", err)
		os.Exit(1)
	}

	app := &cmd.App{
		Config: cfg,
		EC2:    ec2Client,
		TF:     tf.New(cfg.TfDir, false),
		NewSSH: func(host, keyPath string) ssh.Client {
			return ssh.New(host, keyPath)
		},
	}

	if err := cmd.NewRoot(app).Execute(); err != nil {
		os.Exit(1)
	}
}
