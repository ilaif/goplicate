package utils

import (
	"context"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
	"github.com/pkg/errors"
)

type CommandRunner struct {
	Dir string
}

func NewCommandRunner(dir string) *CommandRunner {
	return &CommandRunner{Dir: dir}
}

func (c *CommandRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = c.Dir

	log.Debugf("Running command '%s %s' in directory '%s'", name, strings.Join(args, " "), c.Dir)

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return string(bytes), errors.Wrap(err, "Failed to run command")
	}

	return string(bytes), err
}
