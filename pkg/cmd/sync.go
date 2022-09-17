package cmd

import (
	"github.com/caarlos0/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
	"github.com/ilaif/goplicate/pkg/utils"
)

func NewSyncCmd() *cobra.Command {
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync multiple projects via a configuration file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			_, chToOrigWorkdir, err := utils.ChWorkdir(args)
			if err != nil {
				return err
			}
			defer chToOrigWorkdir()

			config, err := pkg.LoadProjectsConfig()
			if err != nil {
				return err
			}

			for _, project := range config.Projects {
				log.Infof("Syncing project %s...", project.Path)
				log.IncreasePadding()

				if err := utils.Chdir(project.Path); err != nil {
					return err
				}

				projectConfig, err := pkg.LoadProjectConfig()
				if err != nil {
					return err
				}

				if err := pkg.Run(ctx, projectConfig, pkg.NewRunOpts(
					runFlagsOpts.dryRun,
					runFlagsOpts.confirm,
					runFlagsOpts.publish,
					runFlagsOpts.allowDirty,
					runFlagsOpts.force,
					runFlagsOpts.stashChanges,
					runFlagsOpts.baseBranch,
				)); err != nil {
					return errors.Wrapf(err, "Failed to sync project '%s'", project.Path)
				}

				log.DecreasePadding()
				log.Infof("Syncing project %s done", project.Path)
			}

			log.Infof("Syncing complete")

			return nil
		},
	}

	applyRunFlags(syncCmd)

	return syncCmd
}
