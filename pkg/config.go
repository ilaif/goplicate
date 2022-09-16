package pkg

import (
	"github.com/pkg/errors"

	"github.com/ilaif/goplicate/pkg/utils"
)

const (
	defaultProjectConfigFilename  = ".goplicate.yaml"
	defaultProjectsConfigFilename = ".goplicate-projects.yaml"
)

type Target struct {
	Path        string   `yaml:"path"`
	Source      string   `yaml:"source"`
	ParamsPaths []string `yaml:"params"`
}

type Hooks struct {
	Post []string `yaml:"post"`
}

type ProjectConfig struct {
	Targets []Target `yaml:"targets"`
	Hooks   Hooks    `yaml:"hooks"`
}

func LoadProjectConfig() (*ProjectConfig, error) {
	config := &ProjectConfig{}
	if err := utils.ReadYaml(defaultProjectConfigFilename, config); err != nil {
		return nil, errors.Wrap(err, "Failed to load project config")
	}

	return config, nil
}

type ProjectsConfig struct {
	Projects []struct {
		Path string `yaml:"path"`
	} `yaml:"projects"`
}

func LoadProjectsConfig() (*ProjectsConfig, error) {
	config := &ProjectsConfig{}
	if err := utils.ReadYaml(defaultProjectsConfigFilename, config); err != nil {
		return nil, errors.Wrap(err, "Failed to load projects config")
	}

	return config, nil
}
