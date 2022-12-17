package config

import (
	"github.com/pkg/errors"

	"github.com/ilaif/goplicate/pkg/utils"
)

const (
	defaultProjectsConfigFilename = ".goplicate-projects.yaml"
)

func LoadProjectsConfig() (*ProjectsConfig, error) {
	cfg := &ProjectsConfig{}
	if err := utils.ReadYaml(defaultProjectsConfigFilename, cfg); err != nil {
		return nil, errors.Wrap(err, "Failed to load projects config")
	}

	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "Failed to validate projects config")
	}

	return cfg, nil
}

type ProjectsConfig struct {
	Projects []Project `yaml:"projects"`
}

func (pc *ProjectsConfig) Validate() error {
	for _, project := range pc.Projects {
		if err := project.Validate(); err != nil {
			return errors.Wrap(err, "A project is invalid")
		}
	}

	return nil
}

type Project struct {
	Location Source `yaml:"location"`
}

func (p *Project) Validate() error {
	if (p.Location.Repository != "" && p.Location.Path != "") || (p.Location.Repository == "" && p.Location.Path == "") {
		return errors.New("Exactly one of 'repository', 'path' should be specified")
	}

	if err := p.Location.Validate(); err != nil {
		return errors.Wrap(err, "'location' is invalid")
	}

	return nil
}
