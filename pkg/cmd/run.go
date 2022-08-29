package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
)

var runCmdOpts struct {
	dryRun  bool
	confirm bool
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run goplicate on the target repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := pkg.LoadProjectConfig()
		if err != nil {
			return err
		}

		runOpts := pkg.RunOpts{
			DryRun:  runCmdOpts.dryRun,
			Confirm: runCmdOpts.confirm,
		}

		if err := pkg.Run(config, runOpts); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVar(&runCmdOpts.dryRun, "dry-run", false, "do not execute any changes")
	runCmd.Flags().BoolVarP(&runCmdOpts.confirm, "confirm", "y", false, "ask for confirmation")
}
