package config

import (
	"github.com/pkg/errors"

	"github.com/ilaif/goplicate/pkg/utils"
)

const (
	DefaultProjectConfigFilename = ".goplicate.yaml"
)

func LoadProjectConfig() (*ProjectConfig, error) {
	cfg := &ProjectConfig{}
	if err := utils.ReadYaml(DefaultProjectConfigFilename, cfg); err != nil {
		return nil, errors.Wrap(err, "Failed to load project config")
	}

	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "Failed to validate project config")
	}

	return cfg, nil
}

type ProjectConfig struct {
	Targets    []Target `yaml:"targets"`
	Hooks      Hooks    `yaml:"hooks"`
	SyncConfig *Target  `yaml:"sync-config"`
}

func (pc *ProjectConfig) Validate() error {
	for _, target := range pc.Targets {
		if err := target.Validate(); err != nil {
			return errors.Wrap(err, "A target is invalid")
		}
	}

	if pc.SyncConfig != nil {
		if err := pc.SyncConfig.Validate(); err != nil {
			return errors.Wrap(err, "'config-sync' is invalid")
		}
	}

	return nil
}
