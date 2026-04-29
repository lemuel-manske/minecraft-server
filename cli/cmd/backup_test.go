package cmd

import (
	"context"
	"testing"

	"github.com/lkmliz/mc-server/cli/internal/config"
	"github.com/lkmliz/mc-server/cli/internal/ec2"
	"github.com/lkmliz/mc-server/cli/internal/ssh"
	"github.com/lkmliz/mc-server/cli/internal/tf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunBackup_StopsMinecraftFirst(t *testing.T) {
	sshStub, app := makeBackupApp(t)
	require.NoError(t, runBackup(context.Background(), app))
	require.NotEmpty(t, sshStub.RunCmds)
	assert.Equal(t, "sudo systemctl stop minecraft", sshStub.RunCmds[0])
}

func TestRunBackup_StartsMinecraftLast(t *testing.T) {
	sshStub, app := makeBackupApp(t)
	require.NoError(t, runBackup(context.Background(), app))
	assert.Equal(t, "sudo systemctl start minecraft", sshStub.RunCmds[len(sshStub.RunCmds)-1])
}

func TestRunBackup_DownloadsWorld(t *testing.T) {
	sshStub, app := makeBackupApp(t)
	require.NoError(t, runBackup(context.Background(), app))
	require.Len(t, sshStub.DownloadTarPaths, 1)
	assert.Equal(t, "/home/ubuntu/mc/state/world", sshStub.DownloadTarPaths[0][0])
	assert.Contains(t, sshStub.DownloadTarPaths[0][1], "backups/skyblock/")
}

func TestRunBackup_NotRunningErrors(t *testing.T) {
	app := &App{
		Config: &config.Config{InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Inst: &ec2.Instance{ID: "i-abc123", State: "stopped"}},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return &ssh.Stub{} },
	}
	assert.ErrorContains(t, runBackup(context.Background(), app), "server is not running")
}

func TestRunBackup_NoInstanceErrors(t *testing.T) {
	app := &App{
		Config: &config.Config{InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Err: assert.AnError},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return &ssh.Stub{} },
	}
	assert.ErrorContains(t, runBackup(context.Background(), app), "no server found")
}

func makeBackupApp(t *testing.T) (*ssh.Stub, *App) {
	t.Helper()
	sshStub := &ssh.Stub{}
	app := &App{
		Config: &config.Config{
			KeyPath:      "./mc.pem",
			InstanceName: "mc-server",
			Profiles: map[string]config.Profile{
				"skyblock": {WorldDir: "/home/ubuntu/mc/state/world", BackupDir: "./backups/skyblock"},
			},
		},
		EC2:    &ec2.Stub{Inst: &ec2.Instance{ID: "i-abc123", PublicIP: "10.0.0.1", State: "running", Profile: "skyblock"}},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return sshStub },
	}
	return sshStub, app
}
