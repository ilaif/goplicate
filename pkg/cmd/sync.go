package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
	"github.com/ilaif/goplicate/pkg/utils"
)

// var syncCmdOpts struct{}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync multiple projects via a configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		config, err := pkg.LoadProjectsConfig()
		if err != nil {
			return err
		}

		cwd := utils.MustGetwd()
		defer func() {
			_ = os.Chdir(cwd)
		}()

		runOpts := pkg.NewRunOpts(
			runFlagsOpts.dryRun,
			runFlagsOpts.confirm,
			runFlagsOpts.publish,
			runFlagsOpts.allowDirty,
			runFlagsOpts.force,
			runFlagsOpts.baseBranch,
		)

		for _, project := range config.Projects {
			fmt.Printf("Syncing project %s...\n", project.Path)

			if err := utils.Chdir(project.Path); err != nil {
				return err
			}

			projectConfig, err := pkg.LoadProjectConfig()
			if err != nil {
				return err
			}

			if err := pkg.Run(ctx, projectConfig, runOpts); err != nil {
				return err
			}

			fmt.Printf("Syncing project %s done\n", project.Path)
		}

		fmt.Printf("Syncing complete\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	applyRunFlags(syncCmd)
}
