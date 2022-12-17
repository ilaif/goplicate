package config

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// Source a path to a file. Can be from a `repository` if one is specified. Otherwise, assumes a local path.
type Source struct {
	Path string `yaml:"path"`

	Repository RepositoryURI `yaml:"repository"`
	Tag        string        `yaml:"tag"`
	Branch     string        `yaml:"branch"`
	ClonePath  string        `yaml:"clone-path"`
}

func (s *Source) String() string {
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

func (s *Source) Validate() error {
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
