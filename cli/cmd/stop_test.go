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

func TestRunStop_StopsMinecraftService(t *testing.T) {
	sshStub, _, app := makeStopApp(t)
	require.NoError(t, runStop(context.Background(), app))
	assert.Contains(t, sshStub.RunCmds, "sudo systemctl stop minecraft")
}

func TestRunStop_StopsEC2Instance(t *testing.T) {
	_, ec2Stub, app := makeStopApp(t)
	require.NoError(t, runStop(context.Background(), app))
	assert.True(t, ec2Stub.StopCalled)
	assert.Equal(t, "i-abc123", ec2Stub.StopID)
}

func TestRunStop_NotRunningErrors(t *testing.T) {
	app := &App{
		Config: &config.Config{InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Inst: &ec2.Instance{ID: "i-abc123", State: "stopped"}},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return &ssh.Stub{} },
	}
	assert.ErrorContains(t, runStop(context.Background(), app), "server is not running")
}

func TestRunStop_NoInstanceErrors(t *testing.T) {
	app := &App{
		Config: &config.Config{InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Err: assert.AnError},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return &ssh.Stub{} },
	}
	assert.ErrorContains(t, runStop(context.Background(), app), "no server found")
}

func makeStopApp(t *testing.T) (*ssh.Stub, *ec2.Stub, *App) {
	t.Helper()
	sshStub := &ssh.Stub{}
	ec2Stub := &ec2.Stub{Inst: &ec2.Instance{ID: "i-abc123", PublicIP: "10.0.0.1", State: "running", Profile: "skyblock"}}
	app := &App{
		Config: &config.Config{KeyPath: "./mc.pem", InstanceName: "mc-server"},
		EC2:    ec2Stub,
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return sshStub },
	}
	return sshStub, ec2Stub, app
}
