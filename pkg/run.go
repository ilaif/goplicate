package pkg

import (
	"context"

	"github.com/caarlos0/log"
	"github.com/pkg/errors"

	"github.com/ilaif/goplicate/pkg/config"
	"github.com/ilaif/goplicate/pkg/git"
	"github.com/ilaif/goplicate/pkg/utils"
)

type RunOpts struct {
	DryRun       bool
	Confirm      bool
	Publish      bool
	AllowDirty   bool
	Force        bool
	StashChanges bool
	BaseBranch   string
}

func NewRunOpts(
	dryRun, confirm, publish, allowDirty, force, stashChanges bool,
	baseBranch string,
) *RunOpts {
	return &RunOpts{
		DryRun:       dryRun,
		Confirm:      confirm,
		Publish:      publish,
		AllowDirty:   allowDirty,
		Force:        force,
		StashChanges: stashChanges,
		BaseBranch:   baseBranch,
	}
}

func Run(ctx context.Context, cloner git.Cloner, runOpts *RunOpts) error {
	cfg, err := config.LoadProjectConfig()
	if err != nil {
		return err
	}

	updatedTargetPaths := []string{}

	if cfg.SyncConfig != nil {
		target := *cfg.SyncConfig

		if updated, err := RunTarget(ctx, target, cloner, runOpts.DryRun, runOpts.Confirm); err != nil {
			return errors.Wrapf(err, "Target '%s'", target.Path)
		} else if updated {
			updatedTargetPaths = append(updatedTargetPaths, target.Path)
		}

		// Reload the config
		cfg, err = config.LoadProjectConfig()
		if err != nil {
			return err
		}
	}

	publisher := git.NewPublisher(runOpts.BaseBranch, utils.MustGetwd())

	if !runOpts.DryRun && runOpts.Publish {
		if err := publisher.Init(ctx); err != nil {
			return errors.Wrap(err, "Failed to initialize git")
		}

		if !publisher.IsClean() {
			if runOpts.StashChanges {
				restoreStashedChanges, err := publisher.StashChanges(ctx)
				if err != nil {
					return errors.Wrap(err, "Failed to stash changes")
				}

				defer func() {
					if err := restoreStashedChanges(); err != nil {
						log.IncreasePadding()
						log.WithError(err).Warn("Cleanup: Failed to restore stashed changes")
						log.DecreasePadding()
					}
				}()
			} else if !runOpts.AllowDirty {
				return errors.New("Git worktree is not clean. Please commit or stash changes before running again")
			}
		}
	}

	for _, target := range cfg.Targets {
		if updated, err := RunTarget(ctx, target, cloner, runOpts.DryRun, runOpts.Confirm); err != nil {
			return errors.Wrapf(err, "Target '%s'", target.Path)
		} else if updated {
			updatedTargetPaths = append(updatedTargetPaths, target.Path)
		}
	}

	if !runOpts.Force && len(updatedTargetPaths) == 0 {
		return nil
	}

	if !runOpts.DryRun {
		for _, hook := range cfg.Hooks.Post {
			if err := RunHook(ctx, hook); err != nil {
				return err
			}
		}
	}

	if !runOpts.DryRun && runOpts.Publish {
		question := "Do you want to publish the above changes?"
		if answer, err := utils.PromptUserYesNoQuestion(question, runOpts.Confirm); err != nil {
			return err
		} else if answer {
			if err := publisher.Publish(ctx, updatedTargetPaths, runOpts.Confirm); err != nil {
				return errors.Wrap(err, "Failed to publish changes")
			}
		}
	}

	return nil
}
