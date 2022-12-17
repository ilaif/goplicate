package pkg

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	"github.com/ilaif/goplicate/pkg/utils"
)

const (
	defaultProjectConfigFilename  = ".goplicate.yaml"
	defaultProjectsConfigFilename = ".goplicate-projects.yaml"
)

type ProjectConfig struct {
	Targets []Target `yaml:"targets"`
	Hooks   Hooks    `yaml:"hooks"`
}

func (pc *ProjectConfig) Validate() error {
	for _, target := range pc.Targets {
		if err := target.Validate(); err != nil {
			return errors.Wrap(err, "A target is invalid")
		}
	}

	return nil
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

// Target defines a `path` to apply goplicate block snippets on based on the `source` with the supplied `params`
type Target struct {
	Path   string   `yaml:"path"`
	Source Source   `yaml:"source"`
	Params []Source `yaml:"params"`
}

func (t *Target) Validate() error {
	if t.Path == "" {
		return errors.New("'path' cannot be empty")
	}

	if err := t.Source.Validate(); err != nil {
		return errors.Wrap(err, "'source' is invalid")
	}

	for _, param := range t.Params {
		if err := param.Validate(); err != nil {
			return errors.Wrap(err, "A param is invalid")
		}
	}

	return nil
}

// Source a path to a file. Can be from a `repository` if one is specified. Otherwise, assumes a local path.
type Source struct {
	Path string `yaml:"path"`

	Repository RepositoryURI `yaml:"repository"`
	Tag        string        `yaml:"tag"`
	Branch     string        `yaml:"branch"`
	ClonePath  string        `yaml:"clone-path"`
}

func (s Source) String() string {
	if s.Repository == "" {
		return s.Path
	}

	source := string(s.Repository)
	if s.Tag != "" {
		source += fmt.Sprintf("@%s", s.Tag)
	}
	if s.Branch != "" {
		source += fmt.Sprintf("@(%s)", s.Branch)
	}
	if s.Path != "" {
		source += fmt.Sprintf("/%s", s.Path)
	}

	return source
}

func (s Source) Validate() error {
	if s.Repository == "" && s.Path == "" {
		return errors.New("At least one of 'repository', 'path' should be specified")
	}

	if s.Repository != "" {
		if err := s.Repository.Validate(); err != nil {
			return errors.Wrap(err, "'repository' is invalid")
		}
	}

	if s.Tag != "" && s.Branch != "" {
		return errors.New("Only one of 'branch', 'tag' can be specified")
	}

	if s.Repository == "" && (s.Tag != "" || s.Branch != "") {
		return errors.New("'branch' or 'tag' require 'repository' to be specified")
	}

	return nil
}

type RepositoryURI string

func (r RepositoryURI) Validate() error {
	if _, err := url.ParseRequestURI(string(r)); err != nil {
		return errors.Errorf("'%s' is not a valid URI", string(r))
	}

	return nil
}

// Hooks a list of commands to execute after the templating is done.
type Hooks struct {
	Post []string `yaml:"post"`
}

func LoadProjectConfig() (*ProjectConfig, error) {
	config := &ProjectConfig{}
	if err := utils.ReadYaml(defaultProjectConfigFilename, config); err != nil {
		return nil, errors.Wrap(err, "Failed to load project config")
	}

	if err := config.Validate(); err != nil {
		return nil, errors.Wrap(err, "Failed to validate project config")
	}

	return config, nil
}

func LoadProjectsConfig() (*ProjectsConfig, error) {
	config := &ProjectsConfig{}
	if err := utils.ReadYaml(defaultProjectsConfigFilename, config); err != nil {
		return nil, errors.Wrap(err, "Failed to load projects config")
	}

	if err := config.Validate(); err != nil {
		return nil, errors.Wrap(err, "Failed to validate projects config")
	}

	return config, nil
}
