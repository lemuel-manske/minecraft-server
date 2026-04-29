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

func TestRunStatus_QueriesEC2ByInstanceName(t *testing.T) {
	ec2Stub, _, app := makeStatusApp(t)
	require.NoError(t, runStatus(context.Background(), app))
	assert.Equal(t, "mc-server", ec2Stub.GetInstanceName)
}

func TestRunStatus_ChecksServiceStatus(t *testing.T) {
	_, sshStub, app := makeStatusApp(t)
	require.NoError(t, runStatus(context.Background(), app))
	assert.Contains(t, sshStub.OutputCmds, "sudo systemctl status minecraft")
}

func TestRunStatus_NoInstance_ReturnsNoError(t *testing.T) {
	app := &App{
		Config: &config.Config{InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Err: assert.AnError},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return &ssh.Stub{} },
	}
	assert.NoError(t, runStatus(context.Background(), app))
}

func TestRunStatus_StoppedInstance_SkipsSSH(t *testing.T) {
	sshStub := &ssh.Stub{}
	app := &App{
		Config: &config.Config{KeyPath: "./mc.pem", InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Inst: &ec2.Instance{ID: "i-abc123", State: "stopped", Profile: "skyblock"}},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return sshStub },
	}
	require.NoError(t, runStatus(context.Background(), app))
	assert.Empty(t, sshStub.OutputCmds)
}

func makeStatusApp(t *testing.T) (*ec2.Stub, *ssh.Stub, *App) {
	t.Helper()
	ec2Stub := &ec2.Stub{Inst: &ec2.Instance{ID: "i-123", PublicIP: "10.0.0.1", State: "running", Profile: "skyblock"}}
	sshStub := &ssh.Stub{}
	app := &App{
		Config: &config.Config{KeyPath: "./mc.pem", InstanceName: "mc-server"},
		EC2:    ec2Stub,
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return sshStub },
	}
	return ec2Stub, sshStub, app
}
