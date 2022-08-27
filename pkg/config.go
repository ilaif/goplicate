package pkg

const (
	configFilename = ".goplicate.yaml"
)

type RepositoryConfig struct {
	Source  string `yaml:"source"`
	Targets []struct {
		Path   string `yaml:"path"`
		Source string `yaml:"source"`
	} `yaml:"targets"`
}

func LoadRepositoryConfig() (*RepositoryConfig, error) {
	config := &RepositoryConfig{}
	if err := readYaml(configFilename, config); err != nil {
		return nil, err
	}

	return config, nil
}
