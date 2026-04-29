package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lkmliz/mc-server/cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Region(t *testing.T) {
	cfg := loadFixture(t)
	assert.Equal(t, "us-east-1", cfg.Region)
}

func TestLoad_InstanceName(t *testing.T) {
	cfg := loadFixture(t)
	assert.Equal(t, "mc-server", cfg.InstanceName)
}

func TestLoad_OperatorCIDR(t *testing.T) {
	cfg := loadFixture(t)
	assert.Equal(t, "1.2.3.4/32", cfg.OperatorCIDR)
}

func TestLoad_AllowedCIDRs(t *testing.T) {
	cfg := loadFixture(t)
	assert.Equal(t, []string{"5.6.7.8/32", "9.10.11.12/32"}, cfg.AllowedCIDRs)
}

func TestLoad_ProfileModpackDir(t *testing.T) {
	assert.Equal(t, "./modpacks/skyblock", loadProfile(t).ModpackDir)
}

func TestLoad_ProfileBackupKeep(t *testing.T) {
	assert.Equal(t, 5, loadProfile(t).BackupKeep)
}

func TestLoad_ProfileForgeURL(t *testing.T) {
	assert.Equal(t, "https://example.com/forge.jar", loadProfile(t).ForgeURL)
}

func TestLoad_ProfileForgeSHA256(t *testing.T) {
	assert.Equal(t, "abc123", loadProfile(t).ForgeSHA256)
}

func TestGetProfile_Missing(t *testing.T) {
	cfg := &config.Config{Profiles: map[string]config.Profile{}}
	_, err := cfg.GetProfile("nope")
	assert.ErrorContains(t, err, `profile "nope" not found`)
}

func TestLoad_ExpandsEnvVars(t *testing.T) {
	t.Setenv("TEST_CIDR", "9.9.9.9/32")
	f := filepath.Join(t.TempDir(), "mc.yaml")
	require.NoError(t, os.WriteFile(f, []byte(`
region: us-east-1
operator_cidr: "$TEST_CIDR"
`), 0644))
	cfg, err := config.Load(f)
	require.NoError(t, err)
	assert.Equal(t, "9.9.9.9/32", cfg.OperatorCIDR)
}

const fixtureYAML = `
region: us-east-1
key_path: ./mc.pem
tf_dir: ./tf
instance_name: mc-server
operator_cidr: "1.2.3.4/32"
allowed_cidrs:
  - "5.6.7.8/32"
  - "9.10.11.12/32"
profiles:
  skyblock:
    modpack_dir: ./modpacks/skyblock
    backup_dir: ./backups/skyblock
    backup_keep: 5
    forge_url: "https://example.com/forge.jar"
    forge_sha256: "abc123"
`

func loadFixture(t *testing.T) *config.Config {
	t.Helper()
	f := filepath.Join(t.TempDir(), "mc.yaml")
	require.NoError(t, os.WriteFile(f, []byte(fixtureYAML), 0644))
	cfg, err := config.Load(f)
	require.NoError(t, err)
	return cfg
}

func loadProfile(t *testing.T) config.Profile {
	t.Helper()
	cfg := loadFixture(t)
	p, err := cfg.GetProfile("skyblock")
	require.NoError(t, err)
	return p
}
