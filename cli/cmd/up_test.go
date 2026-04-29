package cmd

import (
	"context"
	"strings"
	"testing"

	"github.com/lkmliz/mc-server/cli/internal/config"
	"github.com/lkmliz/mc-server/cli/internal/ec2"
	"github.com/lkmliz/mc-server/cli/internal/ssh"
	"github.com/lkmliz/mc-server/cli/internal/tf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunUp_CallsTerraformApply(t *testing.T) {
	_, tfStub, app := makeUpApp(t)
	require.NoError(t, runUp(context.Background(), app, "skyblock", false))
	assert.True(t, tfStub.ApplyCalled)
}

func TestRunUp_PassesSecurityVarsToTerraform(t *testing.T) {
	_, tfStub, app := makeUpApp(t)
	require.NoError(t, runUp(context.Background(), app, "skyblock", false))
	varStr := strings.Join(tfStub.ApplyVars, " ")
	assert.Contains(t, varStr, `operator_cidr=1.2.3.4/32`)
	assert.Contains(t, varStr, `allowed_cidrs=["5.6.7.8/32","9.10.11.12/32"]`)
	assert.Contains(t, varStr, `forge_url=https://example.com/forge.jar`)
	assert.Contains(t, varStr, `forge_sha256=abc123`)
	assert.Contains(t, varStr, `profile_name=skyblock`)
}

func TestRunUp_UploadsModpack(t *testing.T) {
	sshStub, _, app := makeUpApp(t)
	require.NoError(t, runUp(context.Background(), app, "skyblock", false))
	modpackDir := app.Config.Profiles["skyblock"].ModpackDir
	assert.Equal(t, [][2]string{{modpackDir, "/home/ubuntu/mc"}}, sshStub.UploadTarPaths)
}

func TestRunUp_UploadsServerProperties(t *testing.T) {
	sshStub, _, app := makeUpApp(t)
	require.NoError(t, runUp(context.Background(), app, "skyblock", false))
    modpackDir := app.Config.Profiles["skyblock"].ModpackDir
	want := modpackDir + "/server.properties"
	assert.Equal(t, [][2]string{{want, "/home/ubuntu/mc/server.properties"}}, sshStub.UploadFilePaths)
}

func TestRunUp_StartsMinecraft(t *testing.T) {
	sshStub, _, app := makeUpApp(t)
	require.NoError(t, runUp(context.Background(), app, "skyblock", false))
	assert.Contains(t, sshStub.RunCmds, "sudo systemctl start minecraft")
}

func TestRunUp_UnknownProfileErrors(t *testing.T) {
	_, _, app := makeUpApp(t)
	app.Config.Profiles = map[string]config.Profile{}
	err := runUp(context.Background(), app, "nope", false)
	assert.ErrorContains(t, err, `profile "nope" not found`)
}

func makeUpApp(t *testing.T) (*ssh.Stub, *tf.Stub, *App) {
	t.Helper()
	modpackDir := t.TempDir()
	sshStub := &ssh.Stub{}
	tfStub := &tf.Stub{Outputs: map[string]string{"instance_ip": "10.0.0.1", "instance_id": "i-abc123"}}
	app := &App{
		Config: &config.Config{
			KeyPath:      "./mc.pem",
			InstanceName: "mc-server",
			OperatorCIDR: "1.2.3.4/32",
			AllowedCIDRs: []string{"5.6.7.8/32", "9.10.11.12/32"},
			Profiles: map[string]config.Profile{
				"skyblock": {
					ModpackDir:      modpackDir,
					ForgeURL:        "https://example.com/forge.jar",
					ForgeSHA256:     "abc123",
				},
			},
		},
		EC2:    &ec2.Stub{},
		TF:     tfStub,
		NewSSH: func(host, keyPath string) ssh.Client { return sshStub },
	}
	return sshStub, tfStub, app
}
