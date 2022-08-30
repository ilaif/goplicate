package pkg

import "github.com/ilaif/goplicate/pkg/utils"

const (
	defaultProjectConfigFilename  = ".goplicate.yaml"
	defaultProjectsConfigFilename = ".goplicate-projects.yaml"
)

type Target struct {
	Path       string `yaml:"path"`
	Source     string `yaml:"source"`
	ParamsPath string `yaml:"params"`
}

type ProjectConfig struct {
	Targets []Target `yaml:"targets"`
}

func LoadProjectConfig() (*ProjectConfig, error) {
	config := &ProjectConfig{}
	if err := utils.ReadYaml(defaultProjectConfigFilename, config); err != nil {
		return nil, err
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
		return nil, err
	}

	return config, nil
}
