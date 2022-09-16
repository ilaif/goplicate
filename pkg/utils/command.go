package utils

import (
	"context"
	"os/exec"

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

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "Failed to run command")
	}

	return string(bytes), err
}
