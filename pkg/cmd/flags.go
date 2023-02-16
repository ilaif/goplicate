package cmd

import (
	"github.com/spf13/cobra"
)

var runFlagsOpts struct {
	dryRun         bool
	confirm        bool
	publish        bool
	allowDirty     bool
	force          bool
	stashChanges   bool
	disableCleanup bool
	baseBranch     string
	branch         string
	message        string
}

func applyRunFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&runFlagsOpts.dryRun, "dry-run", false, "do not execute any changes")
	cmd.Flags().BoolVarP(&runFlagsOpts.confirm, "confirm", "y", false, "ask for confirmation")
	cmd.Flags().BoolVar(&runFlagsOpts.publish, "publish", false,
		"publish changes by checking out a new branch, committing, pushing and creating a GitHub pull request",
	)
	cmd.Flags().BoolVar(&runFlagsOpts.allowDirty, "allow-dirty", false, "allow a dirty working tree when publishing")
	cmd.Flags().BoolVar(&runFlagsOpts.force, "force", false, "perform all actions even if there are no updates")
	cmd.Flags().BoolVar(&runFlagsOpts.stashChanges, "stash-changes", false,
		"if the working tree is dirty, stash changes before running, and restore them when done",
	)
	cmd.Flags().BoolVar(&runFlagsOpts.disableCleanup, "disable-cleanup", false, "disable cleanup of cloned repositories")
	cmd.Flags().StringVar(&runFlagsOpts.baseBranch, "base", "", "base git branch to perform updates to")
	cmd.Flags().StringVar(&runFlagsOpts.branch, "branch", "", "name of the new branch to be checked out")
	cmd.Flags().StringVar(&runFlagsOpts.message, "message", "", "pull request description message. supports markdown.")
}
