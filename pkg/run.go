package pkg

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
	"github.com/pkg/errors"
	"github.com/pkg/fileutils"
	"github.com/samber/lo"

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

func Run(ctx context.Context, config *ProjectConfig, cloner git.Cloner, runOpts *RunOpts) error {
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

	updatedTargetPaths := []string{}
	for _, target := range config.Targets {
		if updated, err := runTarget(ctx, target, cloner, runOpts); err != nil {
			return errors.Wrapf(err, "Target '%s'", target.Path)
		} else if updated {
			updatedTargetPaths = append(updatedTargetPaths, target.Path)
		}
	}

	if !runOpts.Force && len(updatedTargetPaths) == 0 {
		return nil
	}

	if !runOpts.DryRun {
		for _, hook := range config.Hooks.Post {
			if err := runHook(ctx, hook); err != nil {
				return err
			}
		}
	}

	if !runOpts.DryRun && runOpts.Publish {
		if !runOpts.Confirm && !utils.AskUserYesNoQuestion("Do you want to publish the above changes?") {
			return errors.New("User aborted")
		}

		if err := publisher.Publish(ctx, updatedTargetPaths, runOpts.Confirm); err != nil {
			return errors.Wrap(err, "Failed to publish changes")
		}
	}

	return nil
}

func runTarget(ctx context.Context, target Target, cloner git.Cloner, runOpts *RunOpts) (bool, error) {
	workdir := utils.MustGetwd()

	sourcePath, err := ResolveSourcePath(ctx, target.Source, workdir, cloner)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to resolve source '%s'", target.Source.String())
	}

	if target.SyncInitial {
		if _, err := os.Stat(target.Path); errors.Is(err, os.ErrNotExist) {
			log.Infof("Syncing initial state of '%s' from '%s'", target.Path, sourcePath)
			if err := fileutils.CopyFile(target.Path, sourcePath); err != nil {
				return false, errors.Wrapf(err, "Failed to copy '%s' to '%s'", sourcePath, target.Path)
			}
		}
	}

	targetBlocks, err := parseBlocksFromFile(target.Path, nil)
	if err != nil {
		return false, errors.Wrap(err, "Failed to parse target blocks")
	}

	params := map[string]interface{}{}
	for _, paramsSource := range target.Params {
		paramsPath, err := ResolveSourcePath(ctx, paramsSource, workdir, cloner)
		if err != nil {
			return false, errors.Wrapf(err, "Failed to resolve source '%s'", paramsSource.String())
		}

		var curParams map[string]interface{}
		if err := utils.ReadYaml(paramsPath, &curParams); err != nil {
			return false, errors.Wrap(err, "Failed to parse params")
		}
		params = lo.Assign(params, curParams)
	}

	sourceBlocks, err := parseBlocksFromFile(sourcePath, params)
	if err != nil {
		return false, errors.Wrap(err, "Failed to parse source blocks")
	}

	anyDiff := false

	for _, targetBlock := range targetBlocks {
		if targetBlock.Name == "" {
			continue
		}

		sourceBlock := sourceBlocks.Get(targetBlock.Name)
		if sourceBlock == nil {
			log.Warnf("Target '%s': Block '%s' not found. Skipping", target.Path, targetBlock.Name)

			continue
		}

		diff := targetBlock.Compare(sourceBlock.Lines)
		if diff != "" {
			log.Infof("Target '%s': Block '%s' needs to be updated. Diff:\n%s\n", target.Path, targetBlock.Name, diff)

			targetBlock.SetLines(sourceBlock.Lines)
			anyDiff = true
		}
	}

	if !anyDiff {
		return false, nil
	}

	if runOpts.DryRun {
		log.Infof("Target '%s': In dry-run mode - Not performing any changes", target.Path)

		return false, nil
	}

	if !runOpts.Confirm && !utils.AskUserYesNoQuestion("Do you want to apply the above changes?") {
		return false, errors.New("User aborted")
	}

	if err := utils.WriteStringToFile(target.Path, targetBlocks.Render()); err != nil {
		return false, err
	}

	log.Infof("Target '%s': Updated", target.Path)

	return true, nil
}

func runHook(ctx context.Context, hook string) error {
	log.Infof("Running post hook '%s'", hook)
	log.IncreasePadding()
	defer log.DecreasePadding()

	cmdParts := strings.Split(hook, " ")
	args := []string{}
	if len(cmdParts) > 0 {
		args = append(args, cmdParts[1:]...)
	}

	outBytes, err := exec.CommandContext(ctx, cmdParts[0], args...).CombinedOutput() // nolint:gosec
	out := string(outBytes)
	if err != nil {
		return errors.Wrapf(err, "Failed to run post hook '%s': %s", hook, out)
	}

	if out != "" {
		log.Infof("Output: %s", out)
	}

	return nil
}
