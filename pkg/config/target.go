package config

import (
	"github.com/pkg/errors"
)

// Target defines a `path` to apply goplicate block snippets on based on the `source` with the supplied `params`
type Target struct {
	Path   string   `yaml:"path"`
	Source Source   `yaml:"source"`
	Params []Source `yaml:"params"`
	// SyncInitial whether to copy the whole file
	// from the source if it doesn't exist.
	SyncInitial bool `yaml:"sync-initial"`
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
