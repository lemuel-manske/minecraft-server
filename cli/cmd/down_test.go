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

func TestRunDown_StopsMinecraft(t *testing.T) {
	sshStub, _, app := makeDownApp(t)
	require.NoError(t, runDown(context.Background(), app))
	assert.Contains(t, sshStub.RunCmds, "sudo systemctl stop minecraft")
}

func TestRunDown_DestroysTerraform(t *testing.T) {
	_, tfStub, app := makeDownApp(t)
	require.NoError(t, runDown(context.Background(), app))
	assert.True(t, tfStub.DestroyCalled)
}

func TestRunDown_NoInstanceErrors(t *testing.T) {
	app := &App{
		Config: &config.Config{InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Err: assert.AnError},
		TF:     &tf.Stub{},
		NewSSH: func(host, keyPath string) ssh.Client { return &ssh.Stub{} },
	}
	assert.ErrorContains(t, runDown(context.Background(), app), "no server found")
}

func makeDownApp(t *testing.T) (*ssh.Stub, *tf.Stub, *App) {
	t.Helper()
	sshStub := &ssh.Stub{}
	tfStub := &tf.Stub{}
	app := &App{
		Config: &config.Config{KeyPath: "./mc.pem", InstanceName: "mc-server"},
		EC2:    &ec2.Stub{Inst: &ec2.Instance{ID: "i-abc123", PublicIP: "10.0.0.1", State: "running", Profile: "skyblock"}},
		TF:     tfStub,
		NewSSH: func(host, keyPath string) ssh.Client { return sshStub },
	}
	return sshStub, tfStub, app
}
