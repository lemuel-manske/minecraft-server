package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	LoaderForge         = "forge"
	LoaderFabric        = "fabric"
	DefaultJavaVersion = 17
)

type LoaderConfig struct {
	Type              string `yaml:"type"`
	MinecraftVersion  string `yaml:"minecraft_version"`
	Version           string `yaml:"version"`
	InstallerURL      string `yaml:"installer_url"`
	InstallerSHA256   string `yaml:"installer_sha256"`
}

type Profile struct {
	JavaVersion         int          `yaml:"java_version"`
	Runtime             LoaderConfig `yaml:"runtime"`
	ModpackDir          string `yaml:"modpack_dir"`
	WorldDir            string `yaml:"world_dir"`
	BackupDir           string `yaml:"backup_dir"`
	BackupKeep          int    `yaml:"backup_keep"`
	Loader              string `yaml:"loader"`
	MinecraftVersion    string `yaml:"minecraft_version"`
	FabricLoaderVersion string `yaml:"fabric_loader_version"`
	ForgeURL            string `yaml:"forge_url"`
	ForgeSHA256         string `yaml:"forge_sha256"`
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

	if p.Runtime.Type != "" {
		p.Loader = p.Runtime.Type
		p.MinecraftVersion = p.Runtime.MinecraftVersion
		p.FabricLoaderVersion = p.Runtime.Version
		p.ForgeURL = p.Runtime.InstallerURL
		p.ForgeSHA256 = p.Runtime.InstallerSHA256
	}
	if p.Loader == "" {
		p.Loader = LoaderForge
	}
	if p.JavaVersion == 0 {
		p.JavaVersion = DefaultJavaVersion
	}
	return p, nil
}
