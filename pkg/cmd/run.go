package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
	"github.com/ilaif/goplicate/pkg/utils"
)

func NewRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Sync the project in the current directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			_, chToOrigWorkdir, err := utils.ChWorkdir(args)
			if err != nil {
				return err
			}
			defer chToOrigWorkdir()

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
