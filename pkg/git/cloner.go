package git

import (
	"context"
	"os"
	"regexp"

	"github.com/caarlos0/log"
	"github.com/pkg/errors"

	"github.com/ilaif/goplicate/pkg/utils"
)

var (
	validPathRegexp = regexp.MustCompile(`[^a-zA-Z0-9_-]`)
)

// Cloner manages cloned git repositories
type Cloner struct {
	repositories map[string]string
}

func NewCloner() *Cloner {
	return &Cloner{
		repositories: make(map[string]string),
	}
}

// Clone clones the repository into a temporary dir and returns it.
// Caches to avoid cloning the same repository twice.
func (c *Cloner) Clone(ctx context.Context, uri string, branch string) (string, error) {
	if tempdir, ok := c.repositories[uri]; ok {
		log.Debugf("Found repository '%s' in cache. Dir is '%s'", uri, tempdir)

		return tempdir, nil
	}

	dirPattern := validPathRegexp.ReplaceAllString("_goplicate_"+uri, "_")
	tempdir, err := os.MkdirTemp(os.TempDir(), dirPattern)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create tempdir '%s'", dirPattern)
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

func (c *Cloner) Close() {
	for uri, tempdir := range c.repositories {
		_ = os.RemoveAll(tempdir)
		delete(c.repositories, uri)
	}
}
