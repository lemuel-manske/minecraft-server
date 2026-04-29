package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Profile struct {
	ModpackDir      string `yaml:"modpack_dir"`
	WorldDir        string `yaml:"world_dir"`
	BackupDir       string `yaml:"backup_dir"`
	BackupKeep      int    `yaml:"backup_keep"`
	ForgeURL        string `yaml:"forge_url"`
	ForgeSHA256     string `yaml:"forge_sha256"`
}

type Config struct {
	Region       string             `yaml:"region"`
	KeyPath      string             `yaml:"key_path"`
	TfDir        string             `yaml:"tf_dir"`
	InstanceName string             `yaml:"instance_name"`
	OperatorCIDR string             `yaml:"operator_cidr"`
	AllowedCIDRs []string           `yaml:"allowed_cidrs"`
	Profiles     map[string]Profile `yaml:"profiles"`
}

func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}

	expanded := os.ExpandEnv(string(raw))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) GetProfile(name string) (Profile, error) {
	p, ok := c.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found", name)
	}
	return p, nil
}
