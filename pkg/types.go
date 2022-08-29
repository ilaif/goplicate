package pkg

import (
	"strings"

	"github.com/pkg/errors"
)

type TargetSource struct {
	Path    string
	Version string
}

func parseTargetSource(s string) (*TargetSource, error) {
	split := strings.Split(s, ":")
	if len(split) > 2 {
		return nil, errors.Errorf("Invalid target source path '%s'. Should be of the form '<path>:<version>'", s)
	}

	targetSource := &TargetSource{
		Path: split[0],
	}

	if len(split) > 1 {
		targetSource.Version = split[1]
	}

	return targetSource, nil
}
