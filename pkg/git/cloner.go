package git

import (
	"context"
	"os"
	"regexp"

	"github.com/caarlos0/log"
	"github.com/pkg/errors"

	"github.com/ilaif/goplicate/pkg/utils"
)

type Cloner interface {
	Clone(
		ctx context.Context,
		uri, branch, fixedClonePath string,
	) (clonePath string, err error)
	Close()
}

var (
	validPathRegexp = regexp.MustCompile(`[^a-zA-Z0-9_-]`)
)

// Cloner manages cloned git repositories
type cloner struct {
	repositories map[string]string
}

func NewCloner() Cloner {
	return &cloner{
		repositories: make(map[string]string),
	}
}

var _ Cloner = &cloner{}

// Clone clones the repository into a temporary dir and returns it.
// Caches to avoid cloning the same repository twice.
func (c *cloner) Clone(
	ctx context.Context,
	uri, branch, fixedClonePath string,
) (string, error) {
	if tempdir, ok := c.repositories[uri]; ok {
		log.Debugf("Found repository '%s' in cache in directory '%s'", uri, tempdir)

		// If there's a clone path and its different from an existing one in
		// the same directory, then we want to symlink to be able to reference it
		if fixedClonePath != "" && tempdir != fixedClonePath {
			if err := os.Symlink(tempdir, fixedClonePath); err != nil {
				return "", errors.Wrapf(err, "Failed to create symlink '%s' for '%s'", tempdir, fixedClonePath)
			}
		}

		return tempdir, nil
	}

	dirPattern := validPathRegexp.ReplaceAllString("_goplicate_"+uri, "_")

	var err error
	tempdir := fixedClonePath
	if tempdir != "" {
		if err := os.MkdirAll(tempdir, 0750); err != nil {
			return "", errors.Wrapf(err, "Failed to create dir '%s'", tempdir)
		}
	} else {
		tempdir, err = os.MkdirTemp(os.TempDir(), dirPattern)
		if err != nil {
			return "", errors.Wrapf(err, "Failed to create tempdir '%s'", dirPattern)
		}
	}

	cmdRunner := utils.NewCommandRunner(tempdir)

	args := []string{"clone", "--depth", "1", uri, "."}
	if branch != "" {
		args = append(args, "--branch", branch)
	}

	log.Infof("Cloning '%s'", uri)

	if output, err := cmdRunner.Run(ctx, "git", args...); err != nil {
		return "", errors.Wrapf(err, "Failed to clone repository '%s': %s", uri, output)
	}

	c.repositories[uri] = tempdir

	return tempdir, nil
}

func (c *cloner) Close() {
	for uri, tempdir := range c.repositories {
		_ = os.RemoveAll(tempdir)
		delete(c.repositories, uri)
	}
}
