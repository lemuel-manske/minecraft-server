package cmd

import (
	"context"
	"io"
	"testing"

	"github.com/lkmliz/mc-server/cli/internal/config"
	"github.com/lkmliz/mc-server/cli/internal/ec2"
	"github.com/lkmliz/mc-server/cli/internal/ssh"
	"github.com/lkmliz/mc-server/cli/internal/tf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunLogs_TailsLogFile(t *testing.T) {
	sshStub, app := makeLogsApp(t)
	require.NoError(t, runLogs(context.Background(), app, io.Discard))
	assert.Contains(t, sshStub.StreamCmds, "tail -f /home/ubuntu/mc/logs/latest.log")
}

func TestRunLogs_NotRunningErrors(t *testing.T) {
	app := &App{
		Config: &config.Config{InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Inst: &ec2.Instance{ID: "i-abc123", State: "stopped"}},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return &ssh.Stub{} },
	}
	assert.ErrorContains(t, runLogs(context.Background(), app, io.Discard), "server is not running")
}

func TestRunLogs_NoInstanceErrors(t *testing.T) {
	app := &App{
		Config: &config.Config{InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Err: assert.AnError},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return &ssh.Stub{} },
	}
	assert.ErrorContains(t, runLogs(context.Background(), app, io.Discard), "no server found")
}

func makeLogsApp(t *testing.T) (*ssh.Stub, *App) {
	t.Helper()
	sshStub := &ssh.Stub{}
	app := &App{
		Config: &config.Config{KeyPath: "./mc.pem", InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Inst: &ec2.Instance{ID: "i-abc123", PublicIP: "10.0.0.1", State: "running", Profile: "skyblock"}},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return sshStub },
	}
	return sshStub, app
}
