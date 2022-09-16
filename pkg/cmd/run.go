package cmd

import (
	"os"

	"github.com/caarlos0/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
	"github.com/ilaif/goplicate/pkg/utils"
)

func newRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Sync the project in the current directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			workDir := lo.Ternary(len(args) > 0, args[0], ".")
			origWorkdir := utils.MustGetwd()
			if err := utils.Chdir(workDir); err != nil {
				return err
			}
			defer func() {
				log.Debugf("Cleanup: Restoring original working directory '%s'", origWorkdir)
				_ = os.Chdir(origWorkdir)
			}()

			config, err := pkg.LoadProjectConfig()
			if err != nil {
				return err
			}

			if err := pkg.Run(ctx, config, pkg.NewRunOpts(
				runFlagsOpts.dryRun,
				runFlagsOpts.confirm,
				runFlagsOpts.publish,
				runFlagsOpts.allowDirty,
				runFlagsOpts.force,
				runFlagsOpts.stashChanges,
				runFlagsOpts.baseBranch,
			)); err != nil {
				return err
			}

			return nil
		},
	}

	applyRunFlags(runCmd)

	return runCmd
}
