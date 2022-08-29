package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
)

var syncCmdOpts struct {
	dryRun  bool
	confirm bool
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync multiple projects via a configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := pkg.LoadProjectsConfig()
		if err != nil {
			return err
		}

		cwd := pkg.MustGetwd()
		defer func() {
			_ = os.Chdir(cwd)
		}()

		runOpts := pkg.RunOpts{
			DryRun:  syncCmdOpts.dryRun,
			Confirm: syncCmdOpts.confirm,
		}

		for _, project := range config.Projects {
			fmt.Printf("Syncing project %s...\n", project.Path)

			if err := pkg.Chdir(project.Path); err != nil {
				return err
			}

			projectConfig, err := pkg.LoadProjectConfig()
			if err != nil {
				return err
			}

			if err := pkg.Run(projectConfig, runOpts); err != nil {
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

	syncCmd.Flags().BoolVar(&syncCmdOpts.dryRun, "dry-run", false, "do not execute any changes")
	syncCmd.Flags().BoolVarP(&syncCmdOpts.confirm, "confirm", "y", false, "ask for confirmation")
}
