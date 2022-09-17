package cmd

import (
	"os"

	"github.com/caarlos0/log"

	"github.com/ilaif/goplicate/pkg/utils"
)

func chWorkdir(args []string) (string, func(), error) {
	workdir := "."
	if len(args) > 0 {
		workdir = args[0]
	}

	origWorkdir := utils.MustGetwd()
	if err := utils.Chdir(workdir); err != nil {
		return workdir, nil, err
	}

	return workdir, func() {
		log.Debugf("Cleanup: Restoring original working directory '%s'", origWorkdir)
		_ = os.Chdir(origWorkdir)
	}, nil
}
