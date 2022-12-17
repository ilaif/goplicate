package cmd

import (
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
	"github.com/ilaif/goplicate/pkg/git"
	"github.com/ilaif/goplicate/pkg/utils"
)

func NewRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Sync the project in the current directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("Executing run command")
			ctx := cmd.Context()

			_, chToOrigWorkdir, err := utils.ChWorkdir(args)
			if err != nil {
				return err
			}
			defer chToOrigWorkdir()

			cloner := git.NewCloner()
			if !runFlagsOpts.disableCleanup {
				defer cloner.Close()
			}

			if err := pkg.Run(ctx, cloner, pkg.NewRunOpts(
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
