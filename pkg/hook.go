package pkg

import (
	"context"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
	"github.com/pkg/errors"
)

func RunHook(ctx context.Context, hook string) error {
	log.Infof("Running post hook '%s'", hook)
	log.IncreasePadding()
	defer log.DecreasePadding()

	cmdParts := strings.Split(hook, " ")
	args := []string{}
	if len(cmdParts) > 0 {
		args = append(args, cmdParts[1:]...)
	}

	outBytes, err := exec.CommandContext(ctx, cmdParts[0], args...).CombinedOutput() // nolint:gosec
	out := string(outBytes)
	if err != nil {
		return errors.Wrapf(err, "Failed to run post hook '%s': %s", hook, out)
	}

	if out != "" {
		log.Infof("Output: %s", out)
	}

	return nil
}
