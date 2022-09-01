package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
)

// var runCmdOpts struct {}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Sync the project in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		config, err := pkg.LoadProjectConfig()
		if err != nil {
			return err
		}

		runOpts := pkg.NewRunOpts(
			runFlagsOpts.dryRun,
			runFlagsOpts.confirm,
			runFlagsOpts.publish,
			runFlagsOpts.allowDirty,
			runFlagsOpts.force,
			runFlagsOpts.baseBranch,
		)
		if err := pkg.Run(ctx, config, runOpts); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	applyRunFlags(runCmd)
}
