package cmd

import (
	"github.com/caarlos0/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
	"github.com/ilaif/goplicate/pkg/git"
	"github.com/ilaif/goplicate/pkg/utils"
)

func NewSyncCmd() *cobra.Command {
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync multiple projects via a configuration file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("Executing sync command")
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

			workdir := utils.MustGetwd()
			cloner := git.NewCloner()
			defer cloner.Close()

			for _, project := range config.Projects {
				projectAbsPath, err := pkg.ResolveSourcePath(ctx, project.Location, workdir, cloner)
				if err != nil {
					return errors.Wrap(err, "Failed to resolve source")
				}

				log.Infof("Syncing project %s...", projectAbsPath)
				log.IncreasePadding()

				if err := utils.Chdir(projectAbsPath); err != nil {
					return err
				}

				projectConfig, err := pkg.LoadProjectConfig()
				if err != nil {
					return err
				}

				if err := pkg.Run(ctx, projectConfig, cloner, pkg.NewRunOpts(
					runFlagsOpts.dryRun,
					runFlagsOpts.confirm,
					runFlagsOpts.publish,
					runFlagsOpts.allowDirty,
					runFlagsOpts.force,
					runFlagsOpts.stashChanges,
					runFlagsOpts.baseBranch,
				)); err != nil {
					return errors.Wrapf(err, "Failed to sync project '%s'", projectAbsPath)
				}

				log.DecreasePadding()
				log.Infof("Syncing project %s done", projectAbsPath)
			}

			log.Infof("Syncing complete")

			return nil
		},
	}

	applyRunFlags(syncCmd)

	return syncCmd
}
