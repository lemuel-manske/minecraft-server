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

func TestRunStart_StartsEC2Instance(t *testing.T) {
	_, ec2Stub, app := makeStartApp(t)
	require.NoError(t, runStart(context.Background(), app))
	assert.True(t, ec2Stub.StartCalled)
	assert.Equal(t, "i-abc123", ec2Stub.StartID)
}

func TestRunStart_WaitsForSSH(t *testing.T) {
	sshStub, _, app := makeStartApp(t)
	require.NoError(t, runStart(context.Background(), app))
	assert.True(t, sshStub.WaitForSSHCalled)
}

func TestRunStart_StartsMinecraftService(t *testing.T) {
	sshStub, _, app := makeStartApp(t)
	require.NoError(t, runStart(context.Background(), app))
	assert.Contains(t, sshStub.RunCmds, "sudo systemctl start minecraft")
}

func TestRunStart_NoInstanceErrors(t *testing.T) {
	app := &App{
		Config: &config.Config{InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Err: assert.AnError},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return &ssh.Stub{} },
	}
	assert.ErrorContains(t, runStart(context.Background(), app), "no server found")
}

func makeStartApp(t *testing.T) (*ssh.Stub, *ec2.Stub, *App) {
	t.Helper()
	sshStub := &ssh.Stub{}
	ec2Stub := &ec2.Stub{
		Inst:   &ec2.Instance{ID: "i-abc123", PublicIP: "10.0.0.1", State: "stopped", Profile: "skyblock"},
		WaitIP: "10.0.0.2",
	}
	app := &App{
		Config: &config.Config{KeyPath: "./mc.pem", InstanceName: "mc-server"},
		EC2:    ec2Stub,
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return sshStub },
	}
	return sshStub, ec2Stub, app
}
