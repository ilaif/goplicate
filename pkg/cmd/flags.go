package cmd

import "github.com/spf13/cobra"

var runFlagsOpts struct {
	dryRun     bool
	confirm    bool
	publish    bool
	allowDirty bool
	force      bool
	baseBranch string
}

func applyRunFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&runFlagsOpts.dryRun, "dry-run", false, "do not execute any changes")
	cmd.Flags().BoolVarP(&runFlagsOpts.confirm, "confirm", "y", false, "ask for confirmation")
	cmd.Flags().BoolVar(&runFlagsOpts.publish, "publish", false,
		"publish changes by checking out a new branch, committing, pushing and creating a GitHub pull request",
	)
	cmd.Flags().BoolVar(&runFlagsOpts.allowDirty, "allow-dirty", false, "allow dirty a working tree when publishing")
	cmd.Flags().BoolVar(&runFlagsOpts.force, "force", false, "perform all actions if there are no changes detected")
	cmd.Flags().StringVar(&runFlagsOpts.baseBranch, "base", "", "the base branch to use when publishing")
}
