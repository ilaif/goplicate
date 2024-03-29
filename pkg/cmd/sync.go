package cmd

import (
	"github.com/caarlos0/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
	"github.com/ilaif/goplicate/pkg/config"
	"github.com/ilaif/goplicate/pkg/git"
	"github.com/ilaif/goplicate/pkg/shared"
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

			cfg, err := config.LoadProjectsConfig()
			if err != nil {
				return err
			}

			workdir := utils.MustGetwd()
			cloner := git.NewCloner()
			if !runFlagsOpts.disableCleanup {
				defer cloner.Close()
			}

			sharedState := &shared.State{
				Message: runFlagsOpts.message,
			}

			for _, project := range cfg.Projects {
				projectAbsPath, err := pkg.ResolveSourcePath(ctx, project.Location, workdir, cloner)
				if err != nil {
					return errors.Wrap(err, "Failed to resolve source")
				}

				log.Infof("Syncing project %s...", projectAbsPath)
				log.IncreasePadding()

				if err := utils.Chdir(projectAbsPath); err != nil {
					return err
				}

				if err := pkg.Run(ctx, cloner, sharedState, pkg.NewRunOpts(
					runFlagsOpts.dryRun,
					runFlagsOpts.confirm,
					runFlagsOpts.publish,
					runFlagsOpts.allowDirty,
					runFlagsOpts.force,
					runFlagsOpts.stashChanges,
					runFlagsOpts.baseBranch,
					runFlagsOpts.branch,
				)); err != nil {
					return errors.Wrapf(err, "Failed to sync project '%s'", projectAbsPath)
				}

				log.DecreasePadding()
				log.Infof("Done syncing project %s", projectAbsPath)
			}

			log.Infof("Syncing complete")

			return nil
		},
	}

	applyRunFlags(syncCmd)

	return syncCmd
}
